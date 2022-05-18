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


#### func (cl \*Client) NewURL(req \*NewURLReq) error
Use NewURL to create a payment link:

```go
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
    req := scanpay.NewURLReq {
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
        Options: &scanpay.Options{
            Headers: map[string]string{
                "X-Cardholder-Ip": "111.222.111.222",
            },
        },
    }
    url, err := client.NewURL(&req)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(url)
}
}
```
### Synchronization
To know when transactions, charges, subscribers and subscriber renewal succeeds, you need to use the synchronization API. It consists of pings which notify you of changes, and the seq request which allows you to pull changes.

#### func (cl \*Client) HandlePing(r \*http.Request) error

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
    APIKey: "1153:YHZIUGQw6NkCIYa3mG6CWcgShnl13xuI7ODFUYuMy0j790Q6ThwBEjxfWFXwJZ0W",
    Host: "api.test.scanpay.dk",
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

#### func (cl \*Client) Seq(req \*scanpay.SeqReq) error

To pull changes since last update, use the Seq() call after receiving a ping.
Store the returned seq-value in a database and use it for the next Seq() call.

```go
package main
import(
    ".."
    "fmt"
)

var client = scanpay.Client{
    APIKey: "1153:YHZIUGQw6NkCIYa3mG6CWcgShnl13xuI7ODFUYuMy0j790Q6ThwBEjxfWFXwJZ0W",
    Host: "api.test.scanpay.dk",
}
var mySeq = uint64(300)

type Acts []scanpay.Act

func seq(pingSeq uint64) {
    for mySeq < pingSeq {
        seqRes, err := client.Seq(&scanpay.SeqReq{ Seq: mySeq })
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

#### func (cl \*Client) Capture(req \*CaptureReq) error
Use Capture to capture a transaction.
```go
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
    req := scanpay.CaptureReq{
        TransactionId: uint64(750),
        Total: "123 DKK",
        Index: 0,
    }
    if err := client.Capture(&req); err != nil {
        fmt.Println("Capture failed:", err)
    } else {
        fmt.Println("Capture succeeded")
    }
}
```
#### func (cl \*Client) Refund(req \*RefundReq) error
Use Refund to refund a captured transaction.
```go
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
    req := scanpay.RefundReq{
        TransactionId: uint64(750),
        Total: "123 DKK",
        Index: 1,
    }
    if err := client.Refund(&req); err != nil {
        fmt.Println("Refund failed:", err)
    } else {
        fmt.Println("Refund succeeded")
    }
}
```
#### func (cl \*Client) Void(req \*VoidReq) error
Use Void to void the amount authorized by the transaction.
```go
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
    APIKey: "1153:YHZIUGQw6NkCIYa3mG6CWcgShnl13xuI7ODFUYuMy0j790Q6ThwBEjxfWFXwJZ0W",
    Host: "api.test.scanpay.dk",
}

func main() {
    req := scanpay.NewURLReq {
        Subscriber: &scanpay.Subscriber{
            Ref: "99",
        },
        Options: &scanpay.Options{
            Headers: map[string]string{
                "X-Cardholder-Ip": "111.222.111.222",
            },
        },
    }
    url, err := client.NewURL(&req)
    if err != nil {
        fmt.Println("NewURL error:", err)
        return
    }
    fmt.Println(url)
}
```

#### func (cl \*Client) Charge(req \*ChargeReq) error
Use Charge to charge a subscriber. The subscriber id is obtained with seq.
```go
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
    req := scanpay.ChargeReq{
        SubscriberId: 30,
        Items: []scanpay.Item{
            {
                Name:"some item",
                Total: "123 DKK",
            },
        },
    }
    if res, err := client.Charge(&req); err != nil {
        fmt.Println("Charge failed:", err)
    } else {
        fmt.Println("Charge succeeded", res)
    }
}
```
#### func (cl \*Client) Renew(req \*RenewReq) error
Use Renew to renew a subscriber, i.e. to attach a new card, if it has expired.
```go
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
```
