package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func Connect() *gorm.DB {
	//dsn := "host=localhost user=blogadmin password=blogadmin dbname=blog port" +
	//	"=5432 sslmode=disable"
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database")
	}
	SeedDevData(db)
	return db
}

func ConnectWithRetry(dsn string, maxAttempts int, delay time.Duration) *gorm.DB {
	var db *gorm.DB
	var err error
	for i := 0; i < maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("successfully connected to database")
			applyMigrations(db)
			SeedDevData(db)
			return db
		}
		log.Printf("Failed to connect to DB (attempt %d/%d): %v", i+1, maxAttempts, err)
		time.Sleep(delay)
	}
	log.Fatalf("Could not connect to database after %d attempts: %v", maxAttempts, err)
	return nil
}
