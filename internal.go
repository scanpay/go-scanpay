package scanpay
import(
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "crypto/subtle"
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "time"
)

type opts struct {
    CardHolderIP string
}

type idempotentResponseErr struct {
    err string
}

func (e *idempotentResponseErr) Error() string {
    return e.err
}

func (c *Client) req(uri string, in interface{}, out interface{}, opts *Options) error {
    var inRdr io.Reader
    reqtype := "GET"
    if in != nil {
        bdata, err := json.Marshal(in)
        if err != nil {
            return err
        }
        inRdr = bytes.NewReader(bdata)
        reqtype = "POST"
    }
    proto := "https://"
    if c.insecure {
        proto = "http://"
    }
    req, err := http.NewRequest(reqtype, proto + c.host + uri, inRdr)
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(c.apikey)) )
    if in != nil {
        req.Header.Set("Content-Type", "application/json; charset=UTF-8")
    }
    if opts != nil && opts.Headers != nil {
        for k, v := range opts.Headers {
            req.Header.Set(k, v)
        }
    }
    idem := req.Header.Get("Idempotency-Key")
    res, err := c.Do(req)
    if err != nil {
        return err
    }
    defer res.Body.Close()
    if idem != "" && res.Header.Get("Idempotency-Status") != "OK" {
        return errors.New("missing idempotency status from response, httpstatus = " + res.Status)
    }
    if res.StatusCode != 200 {
        return &idempotentResponseErr{"scanpay returned " + res.Status}
    }
    if err := json.NewDecoder(res.Body).Decode(out); err != nil {
        switch err.(type) {
        case *json.SyntaxError, *json.UnmarshalTypeError,
            *json.UnsupportedTypeError, *json.UnsupportedValueError:
            return &idempotentResponseErr{"invalid json response from scanpay: " + err.Error()};
        default:
            return err
        }
    }
    return nil
}

/* Debug-only method, may change */
func (c *Client) MakeInsecure() {
    c.insecure = true
}

/* Debug-only method, may change */
func (c *Client) SetHost(host string) {
    c.host = host
}

func (c *Client) signatureIsValid(req *http.Request, body []byte) bool {
    hmacsha2 := hmac.New(sha256.New, []byte(c.apikey))
    hmacsha2.Write(body)
    rawSig := hmacsha2.Sum(nil)
    buf := make([]byte, base64.StdEncoding.EncodedLen(len(rawSig)))
    base64.StdEncoding.Encode(buf, rawSig)
    return subtle.ConstantTimeCompare(buf, []byte(req.Header.Get("X-Signature"))) == 1
}

type LifetimeDuration time.Duration

type internalRenewSubscriberData struct {
    Language   string           `json:"language,omitempty"`
    SuccessURL string           `json:"successurl,omitempty"`
    Lifetime   LifetimeDuration `json:"lifetime,omitempty"`
}

func (d *LifetimeDuration) MarshalText() ([]byte, error) {
    return []byte((*time.Duration)(d).Round(time.Second).String()), nil
}
