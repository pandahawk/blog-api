package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pandahawk/blog-api/internal/database"
	"github.com/pandahawk/blog-api/internal/user"
	"github.com/pandahawk/blog-api/router"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db := database.Connect()

	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatal("failed to create db tables")
	}
	if err := db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 101;").
		Error; err != nil {
		log.Fatal("failed to alter sequence:", err)
	}
	r := gin.Default()
	router.SetupRoutes(r, db)
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
