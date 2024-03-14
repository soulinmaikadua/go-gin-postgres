package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/handlers"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", handlers.CreateUser)
		userGroup.GET("/", handlers.GetUsers)
		userGroup.GET("/:id", handlers.GetUser)
	}
}
