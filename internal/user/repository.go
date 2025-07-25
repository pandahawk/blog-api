package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"gorm.io/gorm"
	"log"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=user
type Repository interface {
	FindAll() ([]*model.User, error)
	FindByID(id uuid.UUID) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	Delete(user *model.User) error
	Update(user *model.User) (*model.User, error)
}

type gormRepository struct {
	db *gorm.DB
}

func (r *gormRepository) FindAll() ([]*model.User, error) {
	var users []*model.User
	err := r.db.Preload("Posts").Find(&users).Error
	return users, err
}

func (r *gormRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Posts").First(&user, id).Error
	return &user, err
}

func (r *gormRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *gormRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *gormRepository) Create(user *model.User) (*model.User, error) {
	err := r.db.Preload("Posts").Create(&user).Error
	return user, err
}

func (r *gormRepository) Update(user *model.User) (*model.User, error) {
	err := r.db.Preload("Posts").Save(&user).Error
	return user, err
}

func (r *gormRepository) Delete(user *model.User) error {
	err := r.db.Delete(&user).Error
	return err
}

func NewDevGormRepository(db *gorm.DB) Repository {

	var user *model.User
	if err := db.First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("no users found... initializing sample data")
		sampleUsers := []*model.User{
			model.NewUser("alice", "alice@example.com"),
			model.NewUser("bob", "bob@example.com"),
			model.NewUser("carl", "carl@example.com"),
			model.NewUser("dave", "dave@example.com"),
			model.NewUser("eve", "eve@example.com"),
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
