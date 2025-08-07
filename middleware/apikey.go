package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func ApiKey() gin.HandlerFunc {
	apiKey := os.Getenv("API_KEY")
	return func(c *gin.Context) {
		if c.GetHeader("X-API-KEY") != apiKey {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		}
		c.Next()
	}
}
