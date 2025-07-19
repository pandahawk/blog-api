package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pandahawk/blog-api/internal/handler"
	"github.com/pandahawk/blog-api/internal/user"
	"gorm.io/gorm"
)

func setupUserRoutes(r *gin.Engine, db *gorm.DB) {
	userRepository := user.NewDevGormRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewHandler(userService)

	v1 := r.Group("/api/v1")
	userGroup := v1.Group("/users")
	userHandler.RegisterRoutes(userGroup)
}

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	if db != nil {
		setupUserRoutes(r, db)
	}

}
