package main

import (
	"github.com/gin-gonic/gin"
	"soniq/internal/server"
)

func main() {
	r := gin.Default()

	// Serve static files (uploaded audio)
	r.Static("/uploads", "./public/uploads")

	// WebSocket endpoint
	r.GET("/ws", server.HandleWebSocket)

	// Audio upload endpoint
	r.POST("/upload", handlers.UploadAudio)

	// Start server
	r.Run(":8080")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}
