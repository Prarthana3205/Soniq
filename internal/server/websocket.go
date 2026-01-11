package server

import (
	"fmt"
	"log"
	"net/http"
	"soniq/internal/server/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var Messages = make(chan string)

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	Clients[conn] = true
	log.Println("New client connected")

	defer func() {
		conn.Close()
		delete(Clients, conn)
	}()

	// Read messages from this client and publish to Redis
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Read error:", err)
				delete(clients, conn)
				break
			}
			redis.PublishMessage(string(message))
		}
	}()

	// Send messages from Redis to this client
	for msg := range Messages {
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				fmt.Println("Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
