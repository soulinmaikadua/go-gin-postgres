package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/models"
	"github.com/soulinmaikadua/go-gin-postgres/internal/routes"
)

func main() {

	// Initialize database
	if err := models.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer models.CloseDB() // Defer closing the database connection until main function exits

	// Initialize Gin router with default middleware
	r := gin.Default()

	// Define a simple ping route for testing server connectivity
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Register application routes for users, posts, and authentication
	routes.AuthRoutes(r) // Register authentication-related routes
	routes.UserRoutes(r) // Register routes related to user management
	routes.PostRoutes(r) // Register routes related to post management

	r.Run(":1234")
}
