package main
import(
    "fmt"
    ".."
)

var client = scanpay.Client{
    APIKey: "1153:YHZIUGQw6NkCIYa3mG6CWcgShnl13xuI7ODFUYuMy0j790Q6ThwBEjxfWFXwJZ0W",
    Host: "api.test.scanpay.dk",
}

func main() {
    trnId := uint64(30)
    data := scanpay.ChargeData{
        Items: []scanpay.Item{
            {
                Name:"some item",
                Total: "123 DKK",
            },
        },
    }
    if res, err := client.Charge(trnId, &data, nil); err != nil {
        fmt.Println("Charge failed:", err)
    } else {
        fmt.Println("Charge succeeded", res)
    }
}
