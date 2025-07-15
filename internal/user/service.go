package user

import (
	"errors"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"strings"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=user

type Service interface {
	GetUser(id int) (User, error)
	CreateUser(req CreateUserRequest) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUser(id int, req UpdateUserRequest) (User, error)
	DeleteUser(id int) error
}

type service struct {
	repo Repository
}

func (s *service) CreateUser(req CreateUserRequest) (User, error) {
	_, found := s.repo.FindByUsername(req.Username)
	if found {
		return User{}, apperrors.NewValidationError("username already exists")
	}
	_, found = s.repo.FindByEmail(req.Email)
	if found {
		return User{}, apperrors.NewValidationError("email already exists")
	}

	user := User{
		Username: req.Username,
		Email:    req.Email,
	}
	return s.repo.Save(user)
}

func (s *service) UpdateUser(id int, req UpdateUserRequest) (User, error) {
	user, found := s.repo.FindByID(id)
	if !found {
		return User{}, apperrors.NewNotFoundError("user", id)
	}

	var validationErrors []string
	if req.Username != nil && strings.TrimSpace(*req.Username) == "" {
		validationErrors = append(validationErrors, "username can not be empty")
	}

	if req.Email != nil && strings.TrimSpace(*req.Email) == "" {
		validationErrors = append(validationErrors, "email cannot be empty")
	}
	if len(validationErrors) > 0 {
		return User{}, apperrors.NewValidationError(validationErrors...)
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	return s.repo.Update(user)
}

func (s *service) DeleteUser(id int) error {
	user, found := s.repo.FindByID(id)
	if !found {
		return apperrors.NewNotFoundError("user", id)
	}

	ok := s.repo.Delete(user)
	if !ok {
		return errors.New("failed to delete user")
	}
	return nil
}

func (s *service) GetUser(id int) (User, error) {
	user, ok := s.repo.FindByID(id)
	if !ok {
		return User{}, apperrors.NewNotFoundError("user", id)
	}
	return user, nil
}

func (s *service) GetAllUsers() ([]User, error) {

	users, ok := s.repo.FindAll()
	if !ok {
		return []User{}, errors.New("failed to get all users")
	}
	return users, nil
}

func NewService(r Repository) Service {
	return &service{repo: r}
}
