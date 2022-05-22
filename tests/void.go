package main
import(
    "fmt"
    "github.com/scanpay/go-scanpay"
)

var client = scanpay.Client{
    APIKey: "1153:YHZIUGQw6NkCIYa3mG6CWcgShnl13xuI7ODFUYuMy0j790Q6ThwBEjxfWFXwJZ0W",
    Host: "api.test.scanpay.dk",
}

func main() {
    req := scanpay.VoidReq{
        TransactionId: uint64(750),
        Index: 0,
    }
    if err := client.Void(&req); err != nil {
        fmt.Println("Void failed:", err)
    } else {
        fmt.Println("Void succeeded")
    }
}
