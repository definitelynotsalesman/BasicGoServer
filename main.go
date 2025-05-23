package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "10000"
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "It works!")
    })

    fmt.Printf("Listening on 0.0.0.0:%s\n", port)
    err := http.ListenAndServe("0.0.0.0:"+port, nil)
    if err != nil {
        panic(err)
    }
}
