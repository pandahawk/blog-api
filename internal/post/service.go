package post

import (
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/shared/model"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=post

type Service interface {
	GetPost(id uuid.UUID) (*model.Post, error)
	CreatePost(req *CreatePostRequest) (*model.Post, error)
	GetPosts() ([]*model.Post, error)
	UpdatePost(id uuid.UUID, req *UpdatePostRequest) (*model.Post, error)
	DeletePost(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func (s service) GetPost(id uuid.UUID) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) CreatePost(req *CreatePostRequest) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetPosts() ([]*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) UpdatePost(id uuid.UUID, req *UpdatePostRequest) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) DeletePost(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
func NewService(repo Repository) Service {
	return &service{repo: repo}
}
