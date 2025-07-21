package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"github.com/pandahawk/blog-api/internal/dto"
	"strings"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=user

type Service interface {
	GetUser(id uuid.UUID) (User, error)
	CreateUser(req dto.CreateUserRequest) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) (User, error)
	DeleteUser(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func (s *service) CreateUser(req dto.CreateUserRequest) (User, error) {
	user, err := s.repo.Create(User{Username: req.Username, Email: req.Email})
	if err != nil {
		if strings.Contains(err.Error(),
			`violates unique constraint "uni_users_username"`) {
			return User{}, apperrors.NewDuplicateError("username")
		}
		if strings.Contains(err.Error(),
			`violates unique constraint "uni_users_email"`) {
			return User{}, apperrors.NewDuplicateError("email")
		}

	}
	return user, nil
}

func (s *service) UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) (User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return User{}, apperrors.NewNotFoundError("user", id)
	}

	if req.Username != nil {
		if strings.TrimSpace(*req.Username) == "" {
			return User{}, apperrors.NewInvalidInputError(
				"username cannot be blank")
		}
		if _, err := s.repo.FindByUsername(*req.Username); err != nil {
			return User{}, apperrors.NewDuplicateError("username")
		}
		user.Username = *req.Username
	}

	if req.Email != nil {
		if strings.TrimSpace(*req.Email) == "" {
			return User{}, apperrors.NewInvalidInputError(
				"email cannot be blank and must be a valid")
		}
		if _, err := s.repo.FindByEmail(*req.Email); err != nil {
			return User{}, apperrors.NewDuplicateError("email")
		}
		user.Email = *req.Email
	}

	return s.repo.Update(user)

	//if req.Username != nil && strings.TrimSpace(*req.Username) == "" {
	//	return User{}, apperrors.NewInvalidInputError("username")
	//}
	//
	//if req.Email != nil && strings.TrimSpace(*req.Email) == "" {
	//	return User{}, apperrors.NewInvalidInputError("email")
	//}
	//
	//if req.Username != nil {
	//	user.Username = *req.Username
	//}
	//if req.Email != nil {
	//	user.Email = *req.Email
	//}

}

func (s *service) DeleteUser(id uuid.UUID) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return apperrors.NewNotFoundError("user", id)
	}

	if err := s.repo.Delete(user); err != nil {
		return errors.New("failed to delete user")
	}
	return nil
}

func (s *service) GetUser(id uuid.UUID) (User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return User{}, apperrors.NewNotFoundError("user", id)
	}
	return user, nil
}

func (s *service) GetAllUsers() ([]User, error) {

	users, err := s.repo.FindAll()
	if err != nil {
		return []User{}, errors.New("failed to get all users")
	}
	return users, nil
}

func NewService(r Repository) Service {
	return &service{repo: r}
}
