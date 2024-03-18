package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	query := "INSERT INTO users (id, name, email, gender, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	userID := uuid.New()
	currentTime := time.Now()
	err := db.QueryRow(query, userID, user.Name, user.Email, user.Gender, user.Password, currentTime, currentTime).Scan(&userID)
	if err != nil {
		fmt.Println("Error creating user:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	defer db.Close()
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

func UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	var user models.UpdateUser
	if err := ctx.BindJSON(&user); err != nil {
		fmt.Println("Error updating user bad request:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Database connection
	db := models.GetConnect()

	// Execute the UPDATE query
	query := "UPDATE users SET "
	var args []interface{}
	argIndex := 1 // Start index for args

	// Check if name is provided
	if user.Name != "" {
		query += "name=$" + strconv.Itoa(argIndex) + ", "
		args = append(args, user.Name)
		argIndex++
	}

	// Check if gender is provided
	if user.Gender != "" {
		query += "gender=$" + strconv.Itoa(argIndex) + ", "
		args = append(args, user.Gender)
		argIndex++
	}

	query += "updated_at=$" + strconv.Itoa(argIndex) + " WHERE id=$" + strconv.Itoa(argIndex+1)
	args = append(args, time.Now(), id)

	_, err := db.Exec(query, args...)
	if err != nil {
		fmt.Println("Error executing query", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(ctx *gin.Context) {
	// Extract user ID from the request parameters
	userID := ctx.Param("id")

	// Database connection
	db := models.GetConnect()

	// Prepare the DELETE query
	query := "DELETE FROM users WHERE id=$1"

	// Execute the DELETE query
	result, err := db.Exec(query, userID)
	if err != nil {
		fmt.Println("SQL error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}
	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking rows affected"})
		return
	}

	if rowsAffected == 0 {
		// No user found with the given ID
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// User deleted successfully
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
