package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pandahawk/blog-api/router"
	"log"
)

func main() {
	r := gin.Default()
	router.SetupRoutes(r)
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
