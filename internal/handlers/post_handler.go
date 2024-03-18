package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/soulinmaikadua/go-gin-postgres/internal/models"
)

func CreatePost(ctx *gin.Context) {
	userID := ctx.GetString("id")
	var post models.Post
	// Parse request body and create user
	if err := ctx.BindJSON(&post); err != nil {
		fmt.Println("Error binding JSON: ", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get database connection
	db := models.GetConnect()

	// Prepare SQL query to insert a new user
	query := "INSERT INTO posts (id, title, details, is_publish, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	postID := uuid.New()
	currentTime := time.Now()
	err := db.QueryRow(query, postID, post.Title, post.Details, post.IsPublish, userID, currentTime, currentTime).Scan(&postID)
	if err != nil {
		fmt.Println("Error creating post:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Set the ID of the created user
	post.ID = postID

	// Respond with the created user
	ctx.JSON(http.StatusCreated, post)
}

func GetPublicPosts(ctx *gin.Context) {
	// Database connection
	db := models.GetConnect()

	// Prepare SQL query to fetch all posts with associated user information
	query := "SELECT p.id, p.title, p.details, p.is_publish, p.created_at, p.updated_at, u.id AS user_id, u.name, u.email FROM posts p JOIN users u ON p.user_id = u.id"

	// Execute the SQL query
	rows, err := db.Query(query)
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

func GetPosts(ctx *gin.Context) {
	// Get user ID from token
	userId := ctx.GetString("id")
	// Database connection
	db := models.GetConnect()

	// Prepare SQL query to fetch all posts with associated user information
	query := "SELECT p.id, p.title, p.details, p.is_publish, p.created_at, p.updated_at, u.id AS user_id, u.name, u.email FROM posts p JOIN users u ON p.user_id = u.id WHERE p.user_id =$1 "

	// Execute the SQL query
	rows, err := db.Query(query, userId)
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

func GetPost(ctx *gin.Context) {
	// Get user ID from token
	userId := ctx.GetString("id")

	// Get post ID from URL parameter
	id := ctx.Param("id")

	// Get a connection to the database
	db := models.GetConnect()

	// Prepare SQL query to fetch a single post by ID
	query := "SELECT p.id, p.title, p.details, p.is_publish, p.created_at, p.updated_at, u.id AS user_id, u.name, u.email FROM posts p JOIN users u ON p.user_id = u.id WHERE p.id=$1 AND p.user_id = $2"

	var post models.ResponsePost
	// Execute query
	err := db.QueryRow(query, id, userId).Scan(&post.ID, &post.Title, &post.Details, &post.IsPublish, &post.CreatedAt, &post.UpdatedAt, &post.User.ID, &post.User.Name, &post.User.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no rows are found
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		// Log the error
		fmt.Println("Error scanning")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the fetched post
	ctx.JSON(http.StatusOK, post)
}

func UpdatePost(ctx *gin.Context) {
	// Extract post ID from the request parameters
	id := ctx.Param("id")

	// Initialize a struct to store the updated post data
	var post models.UpdatePost

	// Bind JSON request body to the post struct
	if err := ctx.BindJSON(&post); err != nil {
		// Return a Bad Request response if JSON binding fails
		fmt.Println("Error updating post bad request:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Database connection
	db := models.GetConnect()

	// Prepare the UPDATE query
	query := "UPDATE posts SET "
	var args []interface{}
	argIndex := 1 // Start index for args

	// Check if title is provided
	if post.Title != "" {
		query += "title=$" + strconv.Itoa(argIndex) + ", "
		args = append(args, post.Title)
		argIndex++
	}

	// Check if details is provided
	if post.Details != "" {
		query += "details=$" + strconv.Itoa(argIndex) + ", "
		args = append(args, post.Details)
		argIndex++
	}

	// Check if is_publish is provided
	if post.IsPublish {
		query += "is_publish=$" + strconv.Itoa(argIndex) + ", "
		args = append(args, post.IsPublish)
		argIndex++
	} else {
		// Set is_publish to false if not provided
		query += "is_publish=false, "
	}

	// Set updated_at field
	query += "updated_at=$" + strconv.Itoa(argIndex) + " WHERE id=$" + strconv.Itoa(argIndex+1)
	args = append(args, time.Now(), id)

	// Execute the UPDATE query
	_, err := db.Exec(query, args...)
	if err != nil {
		// Return an internal server error if query execution fails
		fmt.Println("Error executing query:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating post"})
		return
	}

	// Respond with success message
	ctx.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func DeletePost(ctx *gin.Context) {
	// Extract user ID from the request parameters
	userID := ctx.Param("id")

	// Database connection
	db := models.GetConnect()

	// Prepare the DELETE query
	query := "DELETE FROM posts WHERE id=$1"
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
