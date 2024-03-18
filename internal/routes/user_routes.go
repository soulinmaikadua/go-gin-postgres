package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/handlers"
	"github.com/soulinmaikadua/go-gin-postgres/internal/utils"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		// Public routes
		userGroup.GET("/", handlers.GetUsers)
		userGroup.GET("/:id", handlers.GetUser)
		// Private routes
		// Set verifyToken middleware for the postGroup route group
		userGroup.Use(utils.VerifyToken)
		userGroup.PUT("/:id", handlers.UpdateUser)
		userGroup.DELETE("/:id", handlers.DeleteUser)
	}
}
