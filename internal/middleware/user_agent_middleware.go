package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/soulinmaikadua/go-gin-postgres/internal/utils"
)

func UserAgentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the User-Agent header from the HTTP request
		userAgentString := c.GetHeader("User-Agent")

		// Parse the User-Agent string
		info := utils.ParseUserAgent(userAgentString)

		// Add the User-Agent info to the request context
		c.Set("user_agent", info)

		// Call the next handler function in the chain
		c.Next()
	}
}
