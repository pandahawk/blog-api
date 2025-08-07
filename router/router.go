package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pandahawk/blog-api/internal/post"
	"github.com/pandahawk/blog-api/internal/user"
	"github.com/pandahawk/blog-api/middleware"
	"gorm.io/gorm"
)

func setupResourceRoutes(r *gin.Engine, db *gorm.DB) {
	v1 := r.Group("/api/v1", middleware.ApiKey())

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)
	userGroup := v1.Group("/users")
	userHandler.RegisterRoutes(userGroup)

	postRepository := post.NewRepository(db)
	postService := post.NewService(postRepository)
	postHandler := post.NewHandler(postService)
	postGroup := v1.Group("/posts")
	postHandler.RegisterRoutes(postGroup)

}

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	if db != nil {
		setupResourceRoutes(r, db)
	}

}
