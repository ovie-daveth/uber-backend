package main

import (
	"log"
	"os"
	"time"
	"uber-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	server := gin.New()
	server.Use(gin.Logger())

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins; change to specific origins in production
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Enable if you need to send cookies or authorization headers
		MaxAge:           12 * time.Hour,
	}))

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	routes.AuthRoutes(server)

	server.Run(":" + port)
}
