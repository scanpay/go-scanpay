package scanpay
import(
    "encoding/json"
    "errors"
    "io"
    "io/ioutil"
    "net"
    "net/http"
    "net/url"
    "strconv"
    "time"
)

type Client struct {
    host     string
    apikey   string
    insecure bool
    http.Client
}

func NewClient(apikey string) *Client {
    return &Client{
        host: "api.scanpay.dk",
        apikey: apikey,
        Client: http.Client{
            Transport: &http.Transport{
                Proxy: http.ProxyFromEnvironment,
                TLSHandshakeTimeout: 10 * time.Second,
                Dial: (&net.Dialer{
                    Timeout:   30 * time.Second,
                    KeepAlive: 300 * time.Second,
                }).Dial,
            },
            Timeout: 30 * time.Second,
        },
    }
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

type PaymentURLData struct {
    OrderId     string      `json:"orderid,omitempty"`
    Language    string      `json:"language,omitempty"`
    SuccessURL  string      `json:"successurl,omitempty"`
    AutoCapture bool        `json:"autocapture,omitempty"`
    Items       []Item      `json:"items,omitempty"`
    Subscriber  *Subscriber `json:"subscriber,omitempty"`
    Billing     Billing     `json:"billing,omitempty"`
    Shipping    Shipping    `json:"shipping,omitempty"`
}

type Options struct {
    Headers map[string]string
}

func (c *Client) NewURL(data *PaymentURLData, opts *Options) (string, error) {
    out := struct {
        URL   string `json:"url"`
    }{}
    if err := c.req("/v1/new", data, &out, opts); err != nil {
        return "", err
    }
    if _, err := url.ParseRequestURI(out.URL); err != nil {
        return "", errors.New("Invalid payment URL in new payment url response: " + out.URL)
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

type SeqRes struct {
    Seq     uint64   `json:"seq"`
    Changes []Change `json:"changes"`
}

func (c *Client) Seq(seq uint64, opts *Options) (*SeqRes, error) {
    out := SeqRes{}
    err := c.req("/v1/seq/" + strconv.FormatUint(seq, 10), nil, &out, opts)
    if err != nil {
        return nil, err
    }
    return &out, nil
}

type ChargeData struct {
    OrderId     string     `json:"orderid"`
    AutoCapture bool       `json:"autocapture"`
    Items       []Item     `json:"items"`
    Billing     Billing    `json:"billing"`
    Shipping    Shipping   `json:"shipping"`
}

func (c *Client) Charge(subId uint64, data *ChargeData, opts *Options) error {
    return c.req("/v1/charge/" + strconv.FormatUint(subId, 10), data, nil, opts)
}

/* Check if idempotency-key should be reused */
func IdemReusable(err error) bool {
    _, ok := err.(*idemReusableErr)
    return ok
}
