package main
import(
    ".."
    "fmt"
)

var client = scanpay.Client{
    APIKey: "2881:VwZtDKk+9UuoCWr2xj4G4uFIkQ+Oy4ar/jroMx8NaTE2gmINoehVffYVJTJqfTF2",
    Host: "api.test.scanpay.dk",
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