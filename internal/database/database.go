package database

import (
	"github.com/pandahawk/blog-api/internal/shared/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Connect() *gorm.DB {
	dsn := "host=localhost user=blogadmin password=blogadmin dbname=blog port" +
		"=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	if err := db.AutoMigrate(&model.User{}, &model.Post{}); err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	}

	ApplyCascadeDelete(db)
	return db
}

func ApplyCascadeDelete(db *gorm.DB) {
	sql := `
		ALTER TABLE posts DROP CONSTRAINT IF EXISTS fk_users_posts;
		ALTER TABLE posts
		ADD CONSTRAINT fk_users_posts
		FOREIGN KEY (user_id)
		REFERENCES users(id)
		ON DELETE CASCADE;
	`
	if err := db.Exec(sql).Error; err != nil {
		log.Fatalf("failed to apply ON DELETE CASCADE manually: %v", err)
	}
}
