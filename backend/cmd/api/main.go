// Package main es el punto de entrada de la aplicación.
// Aquí se inicializan las dependencias y se configura el servidor HTTP.
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

	// Configurar modo de Gin
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Crear router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": version,
		})
	})

	// API v1 group (para futuros endpoints)
	api := router.Group("/api/v1")
	{
		// Placeholder para futuros endpoints
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
	}

	// Obtener puerto de variable de entorno o usar default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
