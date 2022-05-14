package main
import(
    "fmt"
    "github.com/scanpay/go-scanpay"
)

var client = scanpay.Client{
    APIKey: "1153:YHZIUGQw6NkCIYa3mG6CWcgShnl13xuI7ODFUYuMy0j790Q6ThwBEjxfWFXwJZ0W",
    Host: "api.test.scanpay.dk"
}

func main() {
    trnId := uint64(750)
    data := scanpay.VoidData{
        Index: 0,
    }
    if err := client.Void(trnId, &data, nil); err != nil {
        fmt.Println("Void failed:", err)
    } else {
        fmt.Println("Void succeeded")
    }
}
