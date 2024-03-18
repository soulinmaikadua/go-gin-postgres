package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/soulinmaikadua/go-gin-postgres/internal/models"
)

var secretKey = []byte("your-secret-key")

func GenerateNewToken(u models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    u.ID,
		"email": u.Email,
		"exp":   time.Now().UTC().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return t, nil
}

func VerifyToken(ctx *gin.Context) {
	// Get the Authorization header
	authHeader := ctx.GetHeader("Authorization")

	// Check if the header is empty or does not contain a bearer token
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		ctx.Abort()
		return
	}

	// Split the header value to extract the token
	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		fmt.Println("print token: ", authToken)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		ctx.Abort()
		return
	}

	// The token is in authToken[1]
	tokenString := authToken[1]

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}

	// Check if the token is valid
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
		ctx.Abort()
		return
	}

	// Token is valid, continue with the next middleware or handler
	fmt.Println("Token is valid")

	// You can also access token claims if needed
	claims, ok := token.Claims.(*jwt.MapClaims)
	if ok {
		ctx.Set("id", (*claims)["id"])
		ctx.Set("email", (*claims)["email"])
	}
	fmt.Println("Token claims:", claims)
	ctx.Next()
}
