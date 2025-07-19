package user

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=user
type Repository interface {
	FindAll() ([]User, error)
	FindByID(id uuid.UUID) (User, error)
	FindByUsername(username string) (User, error)
	FindByEmail(email string) (User, error)
	Create(user User) (User, error)
	Delete(user User) error
	Update(user User) (User, error)
}

type gormRepository struct {
	db *gorm.DB
}

func (r *gormRepository) FindAll() ([]User, error) {
	var users []User
	err := r.db.Preload("Posts").Find(&users).Error
	return users, err
}

func (r *gormRepository) FindByID(id uuid.UUID) (User, error) {
	var user User
	err := r.db.Preload("Posts").First(&user, id).Error
	return user, err
}

func (r *gormRepository) FindByUsername(username string) (User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *gormRepository) FindByEmail(email string) (User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *gormRepository) Create(user User) (User, error) {
	err := r.db.Preload("Posts").Create(&user).Error
	return user, err
}

func (r *gormRepository) Update(user User) (User, error) {
	err := r.db.Preload("Posts").Save(&user).Error
	return user, err
}

func (r *gormRepository) Delete(user User) error {
	err := r.db.Delete(&user).Error
	return err
}

func NewDevGormRepository(db *gorm.DB) Repository {

	var user User
	if err := db.First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("no users found... initializing sample data")
		sampleUsers := []User{
			{Username: "blogger01", Email: "blogger01@example.com"},
			{Username: "blogger02", Email: "blogger02@example.com"},
			{Username: "blogger03", Email: "blogger03@example.com"},
			{Username: "blogger04", Email: "blogger04@example.com"},
			{Username: "blogger05", Email: "blogger05@example.com"},
			{Username: "blogger06", Email: "blogger06@example.com"},
			{Username: "blogger07", Email: "blogger07@example.com"},
			{Username: "blogger08", Email: "blogger08@example.com"},
			{Username: "blogger09", Email: "blogger09@example.com"},
			{Username: "blogger10", Email: "blogger10@example.com"},
		}
		if err := db.Create(&sampleUsers).Error; err != nil {
			log.Fatal("error creating sample users", err)
		}
		log.Println("init sample users successfully")
	}
	return &gormRepository{db: db}
}

func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}
