package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/models"
)

func SearchUsers(ctx *gin.Context) {
	// Get the search query from the request parameters
	searchQuery := ctx.Query("q")
	// Get a connection to the database
	db := models.GetConnect()

	// Construct the query to fetch users
	query := "SELECT id, name, email, gender, created_at, updated_at FROM users WHERE name LIKE '%' || $1 || '%'"

	// Execute the query to fetch users
	rows, err := db.Query(query, searchQuery)
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

func SearchPublicPosts(ctx *gin.Context) {
	// Get the search query from the request parameters
	searchQuery := ctx.Query("q")

	// Database connection
	db := models.GetConnect()

	// Prepare SQL query to fetch posts with associated user information
	query := `
		SELECT p.id, p.title, p.details, p.is_publish, p.created_at, p.updated_at, u.id AS user_id, u.name, u.email
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.title LIKE '%' || $1 || '%'
		OR p.details LIKE '%' || $1 || '%'
		OR u.name LIKE '%' || $1 || '%'
	`

	// Execute the SQL query
	rows, err := db.Query(query, searchQuery)
	if err != nil {
		// Return internal server error if query execution fails
		fmt.Println("Error querying", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Initialize slice to store posts
	var posts []models.ResponsePost

	// Iterate over query results
	for rows.Next() {
		var post models.ResponsePost
		// Scan row into post struct
		err := rows.Scan(&post.ID, &post.Title, &post.Details, &post.IsPublish, &post.CreatedAt, &post.UpdatedAt, &post.User.ID, &post.User.Name, &post.User.Email)
		if err != nil {
			// Return internal server error if scanning fails
			fmt.Println("Error scanning", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Append post to posts slice
		posts = append(posts, post)
	}

	// Check for any errors encountered while iterating over results
	if err := rows.Err(); err != nil {
		// Return internal server error if error encountered
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return posts as JSON response
	ctx.JSON(http.StatusOK, posts)
}
