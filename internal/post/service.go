package post

import (
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/dto"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=post

type Service interface {
	GetPost(id uuid.UUID) (*Post, error)
	CreatePost(req *dto.CreatePostRequest) (*Post, error)
	GetPosts() ([]*Post, error)
	UpdatePost(id uuid.UUID, req *dto.UpdatePostRequest) (*Post, error)
	DeletePost(id uuid.UUID) error
	GetPostSummary(postID uuid.UUID) (*dto.PostSummaryResponse, error)
}

//type UserSummaryService interface {
//	GetUserSummary(userID uuid.UUID) (*dto.UserSummaryResponse, error)
//}

type service struct {
	repo Repository
	//userService UserSummaryService
}

func (s service) GetPostSummary(id uuid.UUID) (*dto.PostSummaryResponse, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &dto.PostSummaryResponse{
		PostID: p.ID,
		Title:  p.Title,
	}, nil
}

func (s service) GetPost(id uuid.UUID) (*Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) CreatePost(req *dto.CreatePostRequest) (*Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetPosts() ([]*Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) UpdatePost(id uuid.UUID, req *dto.UpdatePostRequest) (*Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) DeletePost(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

//	func NewService(repo Repository, us UserSummaryService) Service {
//		return &service{repo: repo, userService: us}
//	}
func NewService(repo Repository) Service {
	return &service{repo: repo}
}
