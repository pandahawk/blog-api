package post

import (
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=post

type Repository interface {
	FindAll() ([]*model.Post, error)
	FindByID(id uuid.UUID) (*model.Post, error)
	Create(user *model.Post) (*model.Post, error)
	Delete(user *model.Post) error
	Update(user *model.Post) (*model.Post, error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) FindAll() ([]*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) FindByID(id uuid.UUID) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Create(user *model.Post) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Delete(user *model.Post) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) Update(user *model.Post) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
