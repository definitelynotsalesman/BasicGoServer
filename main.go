package main

import (
    "fmt"
    "net/http"
    "os"
    "sync"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins
    },
}

var clients = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()

    mutex.Lock()
    clients[conn] = true
    mutex.Unlock()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            break
        }

        // Broadcast to all clients
        mutex.Lock()
        for client := range clients {
            if client != conn {
                client.WriteMessage(websocket.TextMessage, msg)
            }
        }
        mutex.Unlock()
    }

    mutex.Lock()
    delete(clients, conn)
    mutex.Unlock()
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "10000"
    }

    http.HandleFunc("/ws", handleWebSocket)

    fmt.Println("Listening on port", port)
    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        panic("Server error: " + err.Error())
    }
}
