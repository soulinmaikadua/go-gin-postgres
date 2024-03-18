package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/soulinmaikadua/go-gin-postgres/internal/models"
	"github.com/soulinmaikadua/go-gin-postgres/internal/utils"
)

func Login(ctx *gin.Context) {

	var loginInput models.LoginInput
	// Parse request body and create user
	if err := ctx.BindJSON(&loginInput); err != nil {
		fmt.Println("Error binding JSON: ", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Database connection
	db := models.GetConnect()

	// sql query
	query := "SELECT  id, name, email, password, gender, created_at, updated_at FROM users WHERE email=$1"
	var user models.User
	err := db.QueryRow(query, loginInput.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Gender, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		fmt.Println("Error querying user:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Check hashed password
	match := utils.VerifyPassword(loginInput.Password, user.Password)

	if !match {
		// Return status 500 and error message.
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Credentials are not valid"})
		return
	}
	// Generate a new token
	token, err := utils.GenerateNewToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// Set the token in the response header
	ctx.Writer.Header().Set("Authorization", "Bearer "+token)

	// Regenerate the a fake password
	fakePwd, _ := utils.HashPassword("F**K Y_U M@N")
	user.Password = fakePwd

	ctx.JSON(http.StatusOK, user)
}

func Register(ctx *gin.Context) {
	// Parse request body and create user
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		fmt.Println("Error binding JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := uuid.New()
	currentTime := time.Now()
	hashPassword, _ := utils.HashPassword(user.Password)

	// Get database connection
	db := models.GetConnect()

	// Prepare SQL query to insert a new user
	query := "INSERT INTO users (id, name, email, gender, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	err := db.QueryRow(query, userID, user.Name, user.Email, user.Gender, hashPassword, currentTime, currentTime).Scan(&userID)
	if err != nil {
		fmt.Println("Error creating user:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	// Set the ID of the created user
	user.ID = userID

	// Regenerate the a fake password
	fakePwd, _ := utils.HashPassword("F**K Y_U M@N")
	user.Password = fakePwd

	// Respond with the created user
	ctx.JSON(http.StatusCreated, user)
}
