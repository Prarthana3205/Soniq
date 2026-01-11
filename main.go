package main

import (
	"github.com/gin-gonic/gin"
	"soniq/internal/server"
	"soniq/internal/server/handlers"
	"soniq/internal/server/redis"
	"log"
	"net/http"
)

func main() {
	// Initialize Redis connection
	redis.InitRedis()

	// Start listening for messages from Redis and forward them to WebSocket clients
	redis.Subscribe(func(msg string) {
		server.Messages <- msg // Forward to the WebSocket clients
	})

	// Set up Gin router
	r := gin.Default()

	server.StartBroadcastLoop()

	r.Static("/uploads", "./public/uploads")

	r.GET("/ws", server.HandleWebSocket)
	r.POST("/upload", handlers.UploadAudio)

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Start the server with HTTPS
	err := r.RunTLS("0.0.0.0:8080", "server.crt", "server.key") // Update with your actual certificate paths
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
