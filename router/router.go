package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pandahawk/blog-api/internal/user"
)

func SetupRoutes(r *gin.Engine) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	userGroup := r.Group("/users")
	user.RegisterRoutes(userGroup)
}
