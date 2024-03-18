package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/models"
)

func GetUsers(ctx *gin.Context) {
	// Get a connection to the database
	db := models.GetConnect()

	// Construct the query to fetch users
	query := "SELECT id, name, email, gender, created_at, updated_at FROM users"

	// Execute the query to fetch users
	rows, err := db.Query(query)
	if err != nil {
		// Handle database query error
		fmt.Println("Error querying users:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying users"})
		return
	}
	defer rows.Close()

	// Initialize a slice to hold users
	var users []models.ResponseUser

	// Iterate over the result set and scan each row into a ResponseUser struct
	for rows.Next() {
		var user models.ResponseUser
		err := rows.Scan(&user.ID, &user.Name, &user.Gender, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			// Handle scanning error
			fmt.Println("Error scanning user row:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning user row"})
			return
		}
		// Append the user to the slice
		users = append(users, user)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		// Handle iteration error
		fmt.Println("Error iterating over user rows:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over user rows"})
		return
	}

	// Return the list of users as JSON response
	ctx.JSON(http.StatusOK, users)
}

func GetUser(ctx *gin.Context) {
	// Extract the user ID from the request parameters
	id := ctx.Param("id")

	// Get a connection to the database
	db := models.GetConnect()

	// Prepare SQL query to fetch a single user by ID
	query := "SELECT id, name, email, gender, created_at, updated_at FROM users WHERE id = $1"

	// Initialize a ResponseUser struct to hold the fetched user
	var user models.ResponseUser

	// Execute the query and scan the result into the ResponseUser struct
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Gender, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no rows are found
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		// Handle other errors
		fmt.Println("Error querying user:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Respond with the fetched user
	ctx.JSON(http.StatusOK, user)
}

func UpdateUser(ctx *gin.Context) {
	// Extract user ID from the request URL
	id := ctx.Param("id")

	// Bind the JSON request to the UpdateUser struct
	var user models.UpdateUser
	if err := ctx.BindJSON(&user); err != nil {
		// Return a bad request response if JSON binding fails
		fmt.Println("Error updating user: Bad request", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	// Get a connection to the database
	db := models.GetConnect()

	// Construct the UPDATE query
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

	// Set the updated_at field
	query += "updated_at=$" + strconv.Itoa(argIndex) + " WHERE id=$" + strconv.Itoa(argIndex+1)
	args = append(args, time.Now(), id)

	// Execute the UPDATE query
	_, err := db.Exec(query, args...)
	if err != nil {
		// Handle database error
		fmt.Println("Error updating user:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	// Return a success message
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
