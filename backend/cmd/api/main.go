// Package main is the application entry point.
// Here we initialize dependencies and configure the HTTP server.
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort = "8080"
	version     = "0.1.0"
)

func main() {
	log.Println("Gemini Coding Path API - Starting...")

	// Configure Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": version,
		})
	})

	// API v1 group (for future endpoints)
	api := router.Group("/api/v1")
	{
		// Placeholder for future endpoints
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
