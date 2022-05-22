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
