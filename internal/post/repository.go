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
	Create(post *model.Post) (*model.Post, error)
	Delete(post *model.Post) error
	Update(post *model.Post) (*model.Post, error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) FindAll() ([]*model.Post, error) {
	var posts []*model.Post
	err := r.db.Find(&posts).Error
	return posts, err
}

func (r repository) FindByID(id uuid.UUID) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("User").First(&post, id).Error
	return &post, err
}

func (r repository) Create(post *model.Post) (*model.Post, error) {
	err := r.db.Create(post).Error
	return post, err
}

func (r repository) Delete(post *model.Post) error {
	err := r.db.Delete(post).Error
	return err
}

func (r repository) Update(post *model.Post) (*model.Post, error) {
	err := r.db.Save(post).Error
	return post, err
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
