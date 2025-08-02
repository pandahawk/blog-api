package database

import (
	"fmt"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"github.com/pandahawk/blog-api/internal/shared/testdata"
	"gorm.io/gorm"
	"log"
)

func SeedDevData(db *gorm.DB) {

	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		fmt.Println("Skipping seed: users already exist")
		return
	}
	log.Println("Seeding DevData...")

	db.Create(&testdata.SampleUsers)
	db.Create(&testdata.SamplePosts)
}
