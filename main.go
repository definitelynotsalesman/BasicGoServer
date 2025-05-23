func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    fmt.Println("New connection attempt...")
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

        // Broadcast to all other connected clients
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
