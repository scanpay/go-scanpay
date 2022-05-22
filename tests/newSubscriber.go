package main
import(
    ".."
    "fmt"
)

var client = scanpay.Client{
    APIKey: "1153:YHZIUGQw6NkCIYa3mG6CWcgShnl13xuI7ODFUYuMy0j790Q6ThwBEjxfWFXwJZ0W",
    Host: "api.test.scanpay.dk",
}

func main() {
    req := scanpay.NewURLReq {
        Subscriber: &scanpay.Subscriber{
            Ref: "99",
        },
        Options: &scanpay.Options{
            Headers: map[string]string{
                "X-Cardholder-Ip": "111.222.111.222",
            },
        },
    }
    url, err := client.NewURL(&req)
    if err != nil {
        fmt.Println("NewURL error:", err)
        return
    }
    fmt.Println(url)
}
