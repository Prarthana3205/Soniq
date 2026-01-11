package server

import (
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

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	// Register client using the broadcast map
	Clients[conn] = true
	log.Println("WS client connected")

	defer func() {
		conn.Close()
		delete(Clients, conn)
		log.Println("WS client disconnected")
	}()

	// Read messages from this client → Redis
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("WS read error:", err)
				break
			}
			redis.PublishMessage(string(msg))
		}

		conn.Close()
		delete(Clients, conn)
	}()
}
