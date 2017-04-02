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
    Name     string `json:"name"`
    Quantity uint64 `json:"quantity"`
    Price    string `json:"price"`
    SKU      string `json:"sku"`
}

type Billing struct {
    Name    string   `json:"name"`
    Company string   `json:"company"`
    Email   string   `json:"email"`
    Phone   string   `json:"phone"`
    Address []string `json:"address"`
    City    string   `json:"city"`
    Zip     string   `json:"zip"`
    State   string   `json:"state"`
    Country string   `json:"country"`
    VATIN   string   `json:"vatin"`
    GLN     string   `json:"gln"`
}

type Shipping struct {
    Name    string   `json:"name"`
    Company string   `json:"company"`
    Email   string   `json:"email"`
    Phone   string   `json:"phone"`
    Address []string `json:"address"`
    City    string   `json:"city"`
    Zip     string   `json:"zip"`
    State   string   `json:"state"`
    Country string   `json:"country"`
}

type PaymentURLData struct {
    OrderId    string   `json:"orderid"`
    Language   string   `json:"language"`
    SuccessURL string   `json:"successurl"`
    Items      []Item   `json:"items"`
    Billing    Billing  `json:"billing"`
    Shipping   Shipping `json:"shipping"`
}

type Options struct {
    Headers map[string]string
}

func (c *Client) NewURL(data *PaymentURLData, opts *Options) (string, error) {
    out := struct {
        Error string `json:"error"`
        URL   string `json:"url"`
    }{}
    err := c.req("/v1/new", data, &out, opts)
    if err != nil {
        return "", err
    }
    if out.Error != "" {
        return "", errors.New("Error in new payment url response: " + out.Error)
    }
    _, err = url.ParseRequestURI(out.URL)
    if err != nil {
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
type Action struct {
    Action string `json:"act"`
    Time   int64  `json:"time"`
    Total  string `json:"total"`
}

type Change struct {
    Error   string `json:"error"`
    Id      uint64 `json:"id"`
    Rev     uint32 `json:"rev"`
    OrderId string `json:"orderid"`
    Time struct {
        Created    int64 `json:"created"`
        Authorized int64 `json:"authorized"`
    } `json:"time"`
    Actions []Action `json:"acts"`
    Totals struct {
        Authorized string `json:"authorized"`
        Captured   string `json:"captured"`
        Refunded   string `json:"refunded"`
        Left       string `json:"left"`
    }
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
