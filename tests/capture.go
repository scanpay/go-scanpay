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
    trnId := uint64(750)
    data := scanpay.CaptureData{
        Total: "123 DKK",
        Index: 0,
    }
    if err := client.Capture(trnId, &data, nil); err != nil {
        fmt.Println("Capture failed:", err)
    } else {
        fmt.Println("Capture succeeded")
    }
}
