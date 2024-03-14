package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/models"
)

func CreateUser(ctx *gin.Context) {
	// Parse request body and create user
	// Example:
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert user into the database
	// Example:
	// if err := models.CreateUser(&user); err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	ctx.JSON(http.StatusCreated, gin.H{"data": user})
}

func GetUsers(ctx *gin.Context) {

	ctx.JSON(http.StatusCreated, gin.H{"data": "users"})
}
