package database

import (
	"fmt"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func ConnectWithRetry(maxAttempts int, delay time.Duration) *gorm.DB {
	var db *gorm.DB
	var err error

	log.Println("DB_USER:", os.Getenv("DB_USER"))
	log.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	log.Println("DB_NAME:", os.Getenv("DB_NAME"))
	log.Println("DB_HOST:", os.Getenv("DB_HOST"))
	log.Println("DB_PORT:", os.Getenv("DB_PORT"))

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname,
	)

	log.Println("ACTUAL DSN USED:", dsn)

	for i := 0; i < maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("successfully connected to database")
			if err := db.AutoMigrate(&model.User{}, &model.Post{}); err != nil {
				log.Fatalf("AutoMigrate failed: %v", err)
			}
			log.Println("AutoMigrate successful")
			SeedDevData(db)
			return db
		}
		log.Printf("Failed to connect to DB (attempt %d/%d): %v", i+1, maxAttempts, err)
		time.Sleep(delay)
	}
	log.Fatalf("Could not connect to database after %d attempts: %v", maxAttempts, err)
	return nil
}
