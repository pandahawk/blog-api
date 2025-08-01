package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pandahawk/blog-api/internal/post"
	"github.com/pandahawk/blog-api/internal/user"
	"gorm.io/gorm"
)

func setupUserRoutes(r *gin.Engine, db *gorm.DB) {
	v1 := r.Group("/api/v1")

	userRepository := user.NewDevRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)
	userGroup := v1.Group("/users")
	userHandler.RegisterRoutes(userGroup)

	postRepository := post.NewDevRepository(db)
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
		setupUserRoutes(r, db)
	}

}
