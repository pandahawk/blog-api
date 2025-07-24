package post

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=post

type Repository interface {
	FindAll() ([]*Post, error)
	FindByID(id uuid.UUID) (*Post, error)
	Create(user *Post) (*Post, error)
	Delete(user *Post) error
	Update(user *Post) (*Post, error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) FindAll() ([]*Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) FindByID(id uuid.UUID) (*Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Create(user *Post) (*Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) Delete(user *Post) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) Update(user *Post) (*Post, error) {
	//TODO implement me
	panic("implement me")
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
