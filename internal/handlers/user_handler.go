package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/soulinmaikadua/go-gin-postgres/internal/models"
)

func CreateUser(ctx *gin.Context) {
	// Parse request body and create user
	// Example:
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		fmt.Println("Error binding JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get database connection
	db := models.GetConnect()

	// Prepare SQL query to insert a new user
	query := "INSERT INTO users (name, email, gender, password) VALUES ($1, $2, $3, $4) RETURNING id"
	userID := uuid.New()
	fmt.Println(userID.String())
	err := db.QueryRow(query, user.Name, user.Email, user.Gender, user.Password).Scan(&userID)
	if err != nil {
		fmt.Println("Error creating user:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Set the ID of the created user
	user.ID = userID

	// Respond with the created user
	ctx.JSON(http.StatusCreated, user)
}

func GetUsers(ctx *gin.Context) {
	db := models.GetConnect()
	query := "SELECT id, name, email, gender, created_at, updated_at FROM users"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error querying")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.ResponseUser
	for rows.Next() {
		var user models.ResponseUser
		err := rows.Scan(&user.ID, &user.Name, &user.Gender, &user.Email, &user.CreatedAt, &user.UpdatedAt)
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

func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	db := models.GetConnect()
	// Prepare SQL query to fetch a single user by ID
	query := "SELECT id, name, email, gender, created_at, updated_at FROM users WHERE id = $1"

	var user models.ResponseUser
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Gender, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		fmt.Println("Error querying user:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Respond with the fetched user
	ctx.JSON(http.StatusOK, user)
}
