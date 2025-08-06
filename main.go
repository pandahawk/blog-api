package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dsn := os.Getenv("DATABASE_URL")

	fmt.Println("before db connect")
	//db := database.Connect()
	db := database.ConnectWithRetry(dsn, 20, 5*time.Second)
	fmt.Println("after db connect")

	r := gin.Default()

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
