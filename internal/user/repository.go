package user

import (
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"gorm.io/gorm"
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

type repository struct {
	db *gorm.DB
}

func (r *repository) FindAll() ([]*model.User, error) {
	var users []*model.User
	err := r.db.Preload("Posts").Find(&users).Error
	return users, err
}

func (r *repository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Posts").First(&user, id).Error
	return &user, err
}

func (r *repository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *repository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *repository) Create(user *model.User) (*model.User, error) {
	err := r.db.Preload("Posts").Create(&user).Error
	return user, err
}

func (r *repository) Update(user *model.User) (*model.User, error) {
	err := r.db.Preload("Posts").Save(&user).Error
	return user, err
}

func (r *repository) Delete(user *model.User) error {
	err := r.db.Delete(&user).Error
	return err
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
