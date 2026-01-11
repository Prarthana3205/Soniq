package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"soniq/internal/server"
	"soniq/internal/server/handlers"
	"soniq/internal/server/redis"
	"log"
	"net/http"
	"os"
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

	// Serve static uploads
	r.Static("/uploads", "./public/uploads")

	// WebSocket and upload routes
	r.GET("/ws", server.HandleWebSocket)
	r.POST("/upload", handlers.UploadAudio)

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Get port from environment (Railway assigns dynamically)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local testing
	}

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("Listening on %s\n", addr)

	// Start HTTP server (Railway handles HTTPS)
	if err := r.Run(addr); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
