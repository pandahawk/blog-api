package post

import (
	"errors"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"strconv"
	"strings"
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

func validateTitle(title string) error {
	if _, err := strconv.ParseFloat(title, 64); err == nil {
		return apperrors.NewInvalidInputError("title must not be a number")
	}
	if isBlank(title) {
		return apperrors.NewInvalidInputError("title must not be blank")
	}
	if len(title) < 3 {
		return apperrors.NewInvalidInputError("title must have more than 2 characters")
	}
	return nil
}

func isBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func (s service) GetPost(id uuid.UUID) (*model.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("post", id)
	}
	return post, nil
}

func (s service) CreatePost(req *CreatePostRequest) (*model.Post, error) {
	if err := validateTitle(req.Title); err != nil {
		return nil, err
	}

	if isBlank(req.Content) {
		return nil, apperrors.NewInvalidInputError("content must not be blank")
	}

	post := model.NewPost(req.Title, req.Content, req.AuthorID)
	created, err := s.repo.Create(post)
	return created, err
}

func (s service) GetPosts() ([]*model.Post, error) {
	posts, err := s.repo.FindAll()
	if err != nil {
		return nil, errors.New("db error")
	}
	return posts, nil
}

func (s service) UpdatePost(id uuid.UUID, req *UpdatePostRequest) (*model.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("post", id)
	}
	if req.Title != nil {
		err := validateTitle(*req.Title)
		if err != nil {
			return nil, err
		}
		post.Title = *req.Title
	}

	if req.Content != nil {
		if isBlank(*req.Content) {
			return nil, apperrors.NewInvalidInputError("content must not be blank")
		}
		post.Content = *req.Content
	}
	return s.repo.Update(post)
}

func (s service) DeletePost(id uuid.UUID) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return apperrors.NewNotFoundError("post", id)
	}
	err = s.repo.Delete(user)
	if err != nil {
		return errors.New("error deleting post")
	}
	return nil
}
func NewService(repo Repository) Service {
	return &service{repo: repo}
}
