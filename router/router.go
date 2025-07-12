package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pandahawk/blog-api/internal/user"
)

func SetupRoutes(r *gin.Engine) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	userRepository := user.NewDevRepository()
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	userGroup := r.Group("/users")
	userHandler.RegisterRoutes(userGroup)
}
