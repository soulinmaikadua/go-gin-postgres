package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to GO GIN",
		})
	})

	// Define a simple ping route for testing server connectivity
	r.GET("/ping", func(ctx *gin.Context) {
		// Get the User-Agent header from the HTTP request
		userAgent := ctx.GetHeader("User-Agent")
		ipAddress := ctx.ClientIP()
		// Get the User-Agent info from the request context
		info, _ := ctx.Get("user_agent")
		ctx.JSON(200, gin.H{
			"message":    "pong",
			"user_agent": userAgent,
			"info":       info,
			"ip":         ipAddress,
		})
	})

	// Register application routes for users, posts, and authentication
	routes.AuthRoutes(r)   // Register authentication-related routes
	routes.UserRoutes(r)   // Register routes related to user management
	routes.PostRoutes(r)   // Register routes related to post management
	routes.SearchRoutes(r) // Register routes related to search management

	// Create an HTTP server instance
	srv := &http.Server{
		Addr:    ":1234",
		Handler: r,
	}
	// Start the HTTP server in a goroutine
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds
	quit := make(chan os.Signal, 1) // Make the channel buffered with a size of 1

	// Notify the quit channel on SIGINT and SIGTERM signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Wait for signal to arrive
	log.Println("Shutdown Server...")

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	// Wait for the context to be done (timeout or shutdown completed)
	<-ctx.Done()
	log.Println("timeout of 5 seconds.")

	log.Println("Server exiting")
}
