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
        return true // allow all connections
    },
}

var (
    clients = make(map[*websocket.Conn]bool)
    mutex   = sync.Mutex{}
)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    fmt.Println("New connection request...")

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()

    mutex.Lock()
    clients[conn] = true
    mutex.Unlock()

    fmt.Println("Client connected:", conn.RemoteAddr())

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Read error:", err)
            break
        }

        fmt.Println("Received:", string(msg))

        // Broadcast to all other clients
        mutex.Lock()
        for client := range clients {
            if client != conn {
                err := client.WriteMessage(websocket.TextMessage, msg)
                if err != nil {
                    fmt.Println("Write error:", err)
                    client.Close()
                    delete(clients, client)
                }
            }
        }
        mutex.Unlock()
    }

    // Clean up on disconnect
    mutex.Lock()
    delete(clients, conn)
    mutex.Unlock()
    fmt.Println("Client disconnected:", conn.RemoteAddr())
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "WebSocket server is running.")
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "10000"
    }

    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/ws", handleWebSocket)

    addr := "0.0.0.0:" + port
    fmt.Println("Listening on", addr)

    if err := http.ListenAndServe(addr, nil); err != nil {
        fmt.Println("Server error:", err)
        os.Exit(1)
    }
}
