package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/middleware"
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
	r.Use(middleware.UserAgentMiddleware())

	// Define a simple ping route for testing server connectivity
	r.GET("/ping", func(ctx *gin.Context) {
		// Get the User-Agent header from the HTTP request
		userAgent := ctx.GetHeader("User-Agent")
		// Get the User-Agent info from the request context
		info, _ := ctx.Get("user_agent")
		ctx.JSON(200, gin.H{
			"message":    "pong",
			"user_agent": userAgent,
			"info":       info,
		})
	})

	// Register application routes for users, posts, and authentication
	routes.AuthRoutes(r)   // Register authentication-related routes
	routes.UserRoutes(r)   // Register routes related to user management
	routes.PostRoutes(r)   // Register routes related to post management
	routes.SearchRoutes(r) // Register routes related to search management

	r.Run(":1234")
}
