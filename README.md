# Scanpay Go client

The official Go client library for the Scanpay API ([docs](https://docs.scanpay.dk)). You can always e-mail us at [help@scanpay.dk](mailto:help@scanpay.dk), or chat with us on IRC at libera.chat #scanpay

## Installation

```bash
go get github.com/scanpaydk/go-scanpay
```

## Usage

To create a payment link all you need to do is:

```go
package main
import(
    "github.com/scanpaydk/go-scanpay"
    "fmt"
)

func main() {
    client := scanpay.NewClient(" API KEY ")
    data := scanpay.PaymentURLData {
        Items: []scanpay.Item {
            {
                Name:    "Pink Floyd: The Dark Side Of The Moon",
                Quantity: 2,
                Total:   "199.98 DKK",
                SKU:     "fadf23",
            },
        },
    }
    url, err := client.NewURL(&data, nil)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(url)
}
```
