package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/handlers"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.GET("/", handlers.GetUsers)
		userGroup.GET("/:id", handlers.GetUser)
		userGroup.PUT("/:id", handlers.UpdateUser)
		userGroup.DELETE("/:id", handlers.DeleteUser)
	}
}
