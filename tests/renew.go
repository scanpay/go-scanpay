package main
import(
    ".."
    "fmt"
    "time"
)

var client = scanpay.Client{
    APIKey: "1153:YHZIUGQw6NkCIYa3mG6CWcgShnl13xuI7ODFUYuMy0j790Q6ThwBEjxfWFXwJZ0W",
    Host: "api.test.scanpay.dk",
}

func main() {
    req := scanpay.RenewReq {
        SubscriberId: 30,
        Language: "da",
        SuccessURL: "https://scanpay.dk",
        Lifetime: 24 * time.Hour,
    }

    if url, err := client.Renew(&req); err != nil {
        fmt.Println("Renew failed:", err)
    } else {
        fmt.Println("Renew URL:", url)
    }
}
