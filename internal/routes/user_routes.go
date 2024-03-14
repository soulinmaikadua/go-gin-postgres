package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/handlers"
)

func UserRoutes(router *gin.Engine) {
	router.Group("/users")
	router.POST("/", handlers.CreateUser)
	router.GET("/", handlers.GetUsers)
}
