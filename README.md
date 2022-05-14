# Scanpay Go client

The official Go client library for the Scanpay API ([docs](https://docs.scanpay.dk)). You can always e-mail us at [help@scanpay.dk](mailto:help@scanpay.dk), or chat with us on IRC at libera.chat #scanpay

## Installation
```bash
go get github.com/scanpay/go-scanpay
```

## Usage
Create a Scanpay client to start using this library:

```go
var client = scanpay.Client{
    APIKey: " APIKEY ",
}
```

### Payment Link

To create a payment link use NewURL:

```go
package main
import(
    ".."
    "fmt"
)

var client = scanpay.Client{
    APIKey: " APIKEY ",
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
```
### Synchronization
To know when transactions, charges, subscribers and subscriber renewal succeeds, you need to use the synchronization API. It consists of pings which notify you of changes, and the seq request which allows you to pull changes.

#### HandlePing

When changes happen, a **ping** request will be sent to the **ping URL** specified in the Scanpay dashboard.
Use HandlePing to parse the ping request:
```go
package main
import(
    ".."
    "fmt"
    "net/http"
    "os"
)

var client = scanpay.Client{
    APIKey: " APIKEY ",
}

func ping(w http.ResponseWriter, r *http.Request) {
    pingData, err := client.HandlePing(r)
    if err != nil {
        fmt.Println("invalid ping: ", err)
        http.Error(w, "", http.StatusBadRequest)
    } else {
        fmt.Println("Received ping:", pingData)
    }
    os.Exit(0)
}

func main() {
    http.HandleFunc("/ping", ping)
    if err := http.ListenAndServe("localhost:8080", nil); err != nil {
        fmt.Println("http.ListenAndServe failed:", err)
    }
}
```

#### Seq Request

To pull changes since last update, use the Seq() call after receiving a ping.
Store the returned seq-value in a database and use it for the next Seq() call.

```go
package main
import(
    ".."
    "fmt"
)

var client = scanpay.Client{
    APIKey: " APIKEY ",
}
var mySeq = uint64(300)

type Acts []scanpay.Act

func seq(pingSeq uint64) {
    for mySeq < pingSeq {
        seqRes, err := client.Seq(mySeq, nil)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        for _, change := range seqRes.Changes {
            switch change.Type {
            case "transaction", "charge":
                fmt.Printf("Order %s change\n" +
                           "Transaction id: %d\n" +
                           "Revision: %d\n" +
                           "Acts: %#v\n\n",
                           change.OrderId, change.Id, change.Rev, change.Acts)
            case "subscriber":
                fmt.Printf("Subscriber %s change\n" +
                           "Subscriber id: %d\n" +
                           "Revision: %d\n" +
                           "Acts %#v\n\n",
                           change.Ref, change.Id, change.Rev, change.Acts)
            }
        }
        mySeq = seqRes.Seq
        if len(seqRes.Changes) == 0 {
            break
        }
    }
    fmt.Println("final mySeq =", mySeq)
}

func main() {
    pingSeq := uint64(400)
    seq(pingSeq)
}
```
### Transaction Actions

#### Capture
Use Capture to capture a transaction.
```go
package main
import(
    "fmt"
    "github.com/scanpay/go-scanpay"
)

var client = scanpay.Client{
    APIKey: " APIKEY ",
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
```
#### Refund
Use Refund to refund a captured transaction.
```go
package main
import(
    ".."
    "fmt"
)

var client = scanpay.Client{
    APIKey: " APIKEY ",
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
```
#### Void
Use Void to void the amount authorized by the transaction.
```go
package main
import(
    "fmt"
    "github.com/scanpay/go-scanpay"
)

var client = scanpay.Client{ APIKey: " APIKEY " }

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
```
### Subscriptions
Create a subscriber by using NewURL with a Subscriber parameter.
```go
package main
import(
    ".."
    "fmt"
)

var client = scanpay.Client{
    APIKey: " APIKEY ",
}

func main() {
    data := scanpay.PaymentURLData {
        Subscriber: &scanpay.Subscriber{
            Ref: "99",
        },
    }
    opts := scanpay.Options{
        Headers: map[string]string{
            "X-Cardholder-Ip": "111.222.111.222",
        },
    }
    url, err := client.NewURL(&data, &opts)
    if err != nil {
        fmt.Println("NewURL error:", err)
        return
    }
    fmt.Println(url)
}
```

Use Charge to charge a subscriber. The subscriber id is obtained with seq.
```go
package main
import(
    "fmt"
    ".."
)

var client = scanpay.Client{
    APIKey: " APIKEY ",
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
```
Use Renew to renew a subscriber, i.e. to attach a new card, if it has expired.
```go
package main
import(
    ".."
    "fmt"
    "time"
)

var client = scanpay.Client{
    APIKey: " APIKEY ",
}

func main() {
    subId := uint64(30)

    data := scanpay.RenewSubscriberData {
        Language: "da",
        SuccessURL: "https://scanpay.dk",
        Lifetime: 24 * time.Hour,
    }

    if url, err := client.RenewSubscriber(subId, &data, nil); err != nil {
        fmt.Println("Renew failed:", err)
    } else {
        fmt.Println("Renew URL:", url)
    }
}

```
