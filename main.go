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
        return true // allow all origins
    },
}

var clients = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("WebSocket upgrade error:", err)
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

        // Broadcast to all connected clients
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

func handleHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello from Go on Render — WebSocket server is live at /ws")
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "10000"
    }

    http.HandleFunc("/", handleHTTP)
    http.HandleFunc("/ws", handleWebSocket)

    fmt.Printf("Listening on 0.0.0.0:%s\n", port)
    err := http.ListenAndServe("0.0.0.0:"+port, nil)
    if err != nil {
        fmt.Println("Server failed:", err)
        os.Exit(1)
    }
}
