package main

import (
	"fmt"
	"log"
	"os"
	"soniq/internal/server"
	"soniq/internal/server/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Redis
	redis.InitRedis()

	// Start broadcast loop
	server.StartBroadcastLoop()

	// Subscribe Redis → broadcast
	go redis.Subscribe(func(msg string) {
		server.Broadcast <- msg
	})

	// Gin router setup (static, uploads, ws, html)
	r := gin.Default()
	r.Static("/uploads", "./public/uploads")
	r.GET("/ws", server.HandleWebSocket)
	r.POST("/upload", handlers.UploadAudio)
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Println("Listening on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Server failed:", err)
	}
}
