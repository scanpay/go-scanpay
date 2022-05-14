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
    data := scanpay.PaymentURLData {
        OrderId: "a766409",
        Language: "da",
        SuccessURL: "https://insertyoursuccesspage.dk",
        Items: []scanpay.Item {
            {
                Name:    "Pink Floyd: The Dark Side Of The Moon",
                Quantity: 2,
                Total:   "199.98 DKK",
                SKU:     "fadf23",
            },
            {
                Name:    "巨人宏偉的帽子",
                Quantity: 2,
                Total:   "840 DKK",
                SKU:     "124",
            },
        },
        Billing: scanpay.Billing{
            Name:    "John Doe",
            Company: "The Shop A/S",
            Email:   "john@doe.com",
            Phone:   "+4512345678",
            Address: []string{"Langgade 23, 2. th"},
            City:    "Havneby",
            Zip:     "1234",
            State:   "",
            Country: "DK",
            VATIN:   "35413308",
            GLN:     "7495563456235",
        },
        Shipping: scanpay.Shipping{
            Name: "Jan Dåh",
            Company: "The Choppa A/S",
            Email: "jan@doh.com",
            Phone: "+4587654321",
            Address: []string{"Langgade 23, 1. th", "C/O The Choppa"},
            City: "Haveby",
            Zip: "1235",
            State: "",
            Country: "DK",
        },
    }
    opts := scanpay.Options{
        Headers: map[string]string{
            "X-Cardholder-Ip": "111.222.111.222",
        },
    }
    url, err := client.NewURL(&data, &opts)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(url)
}
