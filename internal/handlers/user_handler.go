package handlers

import (
	"fmt"
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
	db := models.GetConnect()
	query := "SELECT * FROM users"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error querying")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Gender, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			fmt.Println("Error scanning")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}
