package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/pandahawk/blog-api/docs"
	"github.com/pandahawk/blog-api/internal/database"
	"github.com/pandahawk/blog-api/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"time"
)

// @title       Blog API
// @version     1.0
// @description This is a simple blog API built with Go and Gin

// @contact.name   Michael Obeng
// @contact.url    https://github.com/pandahawk
// @contact.email  michael@example.com
// @servers [
//   {"url":"http://localhost:8080", "description":"Local"},
//   {"url":"http://89.58.5.201:8080", "description":"Production"}
// ]
// @BasePath  /api/v1

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	_ = godotenv.Load()

	db := database.ConnectWithRetry(5, 5*time.Second)

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.SetupRoutes(r, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
