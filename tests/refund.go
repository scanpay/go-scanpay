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
    trnId := uint64(750)
    data := scanpay.RefundData{
        Total: "123 DKK",
        Index: 1,
    }
    if err := client.Refund(trnId, &data, nil); err != nil {
        fmt.Println("Refund failed:", err)
    } else {
        fmt.Println("Refund succeeded")
    }
}
