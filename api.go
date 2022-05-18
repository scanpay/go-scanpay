package scanpay
import(
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "net"
    "net/http"
    "net/url"
    "strconv"
    "time"
    "unsafe"
)

type Client struct {
    APIKey string
    Host string
    Insecure bool
    HttpClient *http.Client
}

var defaultHttpClient = &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyFromEnvironment,
        TLSHandshakeTimeout: 10 * time.Second,
        Dial: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 300 * time.Second,
        }).Dial,
    },
    Timeout: 30 * time.Second,
}

func (cl *Client) host() string {
    if cl.Host == "" {
        return "api.scanpay.dk"
    }
    return cl.Host
}

func (cl *Client) httpClient() *http.Client {
    if cl.HttpClient == nil {
        return defaultHttpClient
    }
    return cl.HttpClient
}

/* New Payid */
type Item struct {
    Name     string `json:"name,omitempty"`
    Quantity uint64 `json:"quantity,omitempty"`
    Total    string `json:"total,omitempty"`
    SKU      string `json:"sku,omitempty"`
}

type Subscriber struct {
    Ref string `json:"ref,omitempty"`
}

type Billing struct {
    Name    string   `json:"name,omitempty"`
    Company string   `json:"company,omitempty"`
    Email   string   `json:"email,omitempty"`
    Phone   string   `json:"phone,omitempty"`
    Address []string `json:"address,omitempty"`
    City    string   `json:"city,omitempty"`
    Zip     string   `json:"zip,omitempty"`
    State   string   `json:"state,omitempty"`
    Country string   `json:"country,omitempty"`
    VATIN   string   `json:"vatin,omitempty"`
    GLN     string   `json:"gln,omitempty"`
}

type Shipping struct {
    Name    string   `json:"name,omitempty"`
    Company string   `json:"company,omitempty"`
    Email   string   `json:"email,omitempty"`
    Phone   string   `json:"phone,omitempty"`
    Address []string `json:"address,omitempty"`
    City    string   `json:"city,omitempty"`
    Zip     string   `json:"zip,omitempty"`
    State   string   `json:"state,omitempty"`
    Country string   `json:"country,omitempty"`
}

type NewURLReq struct {
    OrderId     string      `json:"orderid,omitempty"`
    Language    string      `json:"language,omitempty"`
    SuccessURL  string      `json:"successurl,omitempty"`
    AutoCapture bool        `json:"autocapture,omitempty"`
    Items       []Item      `json:"items,omitempty"`
    Subscriber  *Subscriber `json:"subscriber,omitempty"`
    Billing     Billing     `json:"billing,omitempty"`
    Shipping    Shipping    `json:"shipping,omitempty"`
    Options     *Options    `json:"-"`
}

type Options struct {
    Headers map[string]string
}

func (c *Client) NewURL(data *NewURLReq) (string, error) {
    out := struct {
        URL   string `json:"url"`
    }{}
    if err := c.req("/v1/new", data, &out, data.Options); err != nil {
        return "", err
    }
    if _, err := url.ParseRequestURI(out.URL); err != nil {
        return "", fmt.Errorf("Invalid payment URL in new payment url response: %s", out.URL)
    }
    return out.URL, nil
}

type Ping struct {
    ShopId uint64 `json:"shopid"`
    Seq    uint64 `json:"seq"`
}

func (c *Client) HandlePing(req *http.Request) (*Ping, error) {
    body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1024))
    if err != nil {
        return nil, err
    }
    if !c.signatureIsValid(req, body) {
        return nil, errors.New("Invalid ping signature")
    }
    ping := Ping{}
    if err := json.Unmarshal(body, &ping); err != nil {
        return nil, err
    }
    if ping.ShopId == 0 {
        return nil, errors.New("Missing field")
    }
    return &ping, nil
}

/* Sequence request */
type Act struct {
    Action string `json:"act"`
    Time   int64  `json:"time"`
    Total  string `json:"total"`
}

type Change struct {
    Type    string `json:"type"`
    Error   string `json:"error"`
    Id      uint64 `json:"id"`
    Rev     uint32 `json:"rev"`
    OrderId string `json:"orderid"`
    Time struct {
        Created    int64 `json:"created"`
        Authorized int64 `json:"authorized"`
    } `json:"time"`
    Acts []Act `json:"acts"`
    Totals struct {
        Authorized string `json:"authorized"`
        Captured   string `json:"captured"`
        Refunded   string `json:"refunded"`
        Left       string `json:"left"`
    } `json:"totals"`

    /* type=subscriber specific data */
    Ref     string `json:"ref"`

    /* type=charge specific data */
    Subscriber struct {
        Id  uint64 `json:"id"`
        Ref string `json:"ref"`
    } `json:"subscriber"`
}

type SeqReq struct {
    Seq     uint64   `json:"seq"`
    Options *Options `json:"-"`
}

type SeqRes struct {
    Seq     uint64   `json:"seq"`
    Changes []Change `json:"changes"`
}

func (c *Client) Seq(data *SeqReq) (*SeqRes, error) {
    out := SeqRes{}
    err := c.req("/v1/seq/" + strconv.FormatUint(data.Seq, 10), nil, &out, data.Options)
    if err != nil {
        return nil, err
    }
    return &out, nil
}

type ChargeReq struct {
    SubscriberId uint64     `json:"-"`
    OrderId      string     `json:"orderid"`
    AutoCapture  bool       `json:"autocapture"`
    Items        []Item     `json:"items"`
    Billing      Billing    `json:"billing"`
    Shipping     Shipping   `json:"shipping"`
    Options      *Options   `json:"-"`
}

type ChargeRes struct {
    Id uint64 `json:"id"`
    Totals struct {
        Authorized string `json:"authorized"`
    } `json:"totals"`
}

func (c *Client) Charge(data *ChargeReq) (*ChargeRes, error) {
    out := ChargeRes{}
    err := c.req("/v1/subscribers/" + strconv.FormatUint(data.SubscriberId, 10) + "/charge",
                 data, &out, data.Options)
    if err != nil {
        return nil, err
    }
    return &out, nil
}

type ActReq struct {
    TransactionId uint64   `json:"-"`
    Total   string   `json:"total"`
    Index   uint64   `json:"index"`
    Options *Options `json:"-"`
}

type CaptureReq ActReq
type RefundReq ActReq
type VoidReq ActReq

func (c *Client) Capture(data *CaptureReq) error {
    return c.req("/v1/transactions/" + strconv.FormatUint(data.TransactionId, 10) + "/capture",
                 data, nil, data.Options)
}

func (c *Client) Refund(data *RefundReq) error {
    return c.req("/v1/transactions/" + strconv.FormatUint(data.TransactionId, 10) + "/refund",
                 data, nil, data.Options)
}

func (c *Client) Void(data *VoidReq) error {
    return c.req("/v1/transactions/" + strconv.FormatUint(data.TransactionId, 10) + "/void",
                 data, nil, data.Options)
}

/* Check if idempotency-key should be reused */
func IsIdempotentResponseError(err error) bool {
    _, ok := err.(*idempotentResponseErr)
    return ok
}

type RenewReq struct {
    SubscriberId uint64        `json:"-"`
    Language     string        `json:"language,omitempty"`
    SuccessURL   string        `json:"successurl,omitempty"`
    Lifetime     time.Duration `json:"lifetime,omitempty"`
    Options      *Options      `json:"-"`
}

func (c *Client) Renew(data *RenewReq) (string, error) {
    out := struct {
        URL   string `json:"url"`
    }{}
    /* Convert to internal representation to allow for alternate duraton formatting */
    d := (*internalRenewSubscriberData)(unsafe.Pointer(data))
    if err := c.req(fmt.Sprintf("/v1/subscribers/%d/renew", data.SubscriberId), d, &out, data.Options); err != nil {
        return "", err
    }
    if _, err := url.ParseRequestURI(out.URL); err != nil {
        return "", fmt.Errorf("Invalid renew URL in new renew subscriber response: %s", out.URL)
    }
    return out.URL, nil
}
