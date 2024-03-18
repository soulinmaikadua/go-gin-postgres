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

func CreatePost(ctx *gin.Context) {
	var post models.Post
	// Parse request body and create user
	if err := ctx.BindJSON(&post); err != nil {
		fmt.Println("Error binding JSON: ", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// This will be changed in the future
	userID := "21d936b2-fc68-4dbc-8f20-c7c88a7e4a39"

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

func GetPosts(ctx *gin.Context) {

	db := models.GetConnect()
	query := "SELECT p.id, p.title, p.details, p.is_publish, p.created_at, p.updated_at, u.id AS user_id, u.name, u.email FROM posts p JOIN users u ON p.user_id = u.id "
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error querying")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var posts []models.ResponsePost
	for rows.Next() {
		var post models.ResponsePost
		err := rows.Scan(&post.ID, &post.Title, &post.Details, &post.IsPublish, &post.CreatedAt, &post.UpdatedAt, &post.User.ID, &post.User.Name, &post.User.Email)
		if err != nil {
			fmt.Println("Error scanning")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func GetPost(ctx *gin.Context) {
	id := ctx.Param("id")
	db := models.GetConnect()
	query := "SELECT p.id, p.title, p.details, p.is_publish, p.created_at, p.updated_at, u.id AS user_id, u.name, u.email FROM posts p JOIN users u ON p.user_id = u.id WHERE p.id=$1"

	var post models.ResponsePost
	err := db.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Details, &post.IsPublish, &post.CreatedAt, &post.UpdatedAt, &post.User.ID, &post.User.Name, &post.User.Email)
	if err != nil {
		fmt.Println("Error scanning")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, post)
}

func UpdatePost(ctx *gin.Context) {
	id := ctx.Param("id")

	var post models.UpdatePost
	// TODO: json binding
	if err := ctx.BindJSON(&post); err != nil {
		fmt.Println("Error updating user bad request:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Database connection
	db := models.GetConnect()

	// Execute the UPDATE query
	query := "UPDATE posts SET "
	var args []interface{}
	argIndex := 1 // Start index for args

	// Check if title is provided
	if post.Title != "" {
		query += "title=$" + strconv.Itoa(argIndex) + ", "
		args = append(args, post.Title)
		argIndex++
	}

	// Check if detail is provided
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
		query += "is_publish=false"
	}
	// Default once
	query += "updated_at=$" + strconv.Itoa(argIndex) + " WHERE id=$" + strconv.Itoa(argIndex+1)
	args = append(args, time.Now(), id)

	_, err := db.Exec(query, args...)
	if err != nil {
		fmt.Println("error executing query", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating post"})
		return
	}
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
