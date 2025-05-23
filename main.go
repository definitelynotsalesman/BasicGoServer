package main

import (
    "fmt"
    "net/http"
    "os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        fmt.Println("ERROR: PORT not set — Render needs this")
        port = "10000" // for local testing fallback
    } else {
        fmt.Println("Using PORT:", port)
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Render check success")
    })

    addr := "0.0.0.0:" + port
    fmt.Println("Listening on", addr)

    if err := http.ListenAndServe(addr, nil); err != nil {
        fmt.Println("ListenAndServe error:", err)
        os.Exit(1)
    }
}
