package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/handlers"
	"github.com/soulinmaikadua/go-gin-postgres/internal/utils"
)

func PostRoutes(router *gin.Engine) {
	postGroup := router.Group("/posts")
	{
		// Public routes
		postGroup.GET("/", handlers.GetPublicPosts)
		// Set verifyToken middleware for the postGroup route group
		postGroup.Use(utils.VerifyToken)

		postGroup.POST("/", handlers.CreatePost)
		postGroup.GET("/my-posts", handlers.GetPosts)
		postGroup.GET("/:id", handlers.GetPost)
		postGroup.PUT("/:id", handlers.UpdatePost)
		postGroup.DELETE("/:id", handlers.DeletePost)
	}
}
