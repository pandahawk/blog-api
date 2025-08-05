package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
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
