package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/handlers"
)

func SearchRoutes(router *gin.Engine) {
	userGroup := router.Group("/search")
	{
		// Public routes
		userGroup.GET("/users", handlers.SearchUsers)
		userGroup.GET("/posts", handlers.SearchPublicPosts)

	}
}
