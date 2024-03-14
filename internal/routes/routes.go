package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/routes"
)

func SetupRoutes(router *gin.Engine) {
	routes.UserRoutes(router)
}
