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
	defer models.CloseDB()

	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// all routes
	routes.UserRoutes(r)
	routes.PostRoutes(r)
	routes.AuthRoutes(r)

	r.Run(":1234")
}
