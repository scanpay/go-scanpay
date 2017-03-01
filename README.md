# Scanpay Go client

This is a Go client library for Scanpay. You can find the documentation at [docs.scanpay.dk](https://docs.scanpay.dk/).

## Installation

```bash
go get github.com/scanpaydk/go-scanpay
```

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
                Price:   "99.99 DKK",
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
