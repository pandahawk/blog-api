package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"strings"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=user

type Service interface {
	GetUser(id uuid.UUID) (User, error)
	CreateUser(req CreateUserRequest) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUser(id uuid.UUID, req UpdateUserRequest) (User, error)
	DeleteUser(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func (s *service) CreateUser(req CreateUserRequest) (User, error) {
	var validationErrors []string
	if _, err := s.repo.FindByUsername(req.Username); err == nil {
		validationErrors = append(validationErrors, "username already exists")
	}

	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		validationErrors = append(validationErrors, "email already exists")
	}

	if len(validationErrors) > 0 {
		return User{}, apperrors.NewValidationError(validationErrors...)
	}

	user, err := s.repo.Create(User{Username: req.Username, Email: req.Email})
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *service) UpdateUser(id uuid.UUID, req UpdateUserRequest) (User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return User{}, apperrors.NewNotFoundError("user", id)
	}

	var validationErrors []string
	if req.Username != nil && strings.TrimSpace(*req.Username) == "" {
		validationErrors = append(validationErrors, "username can not be empty")
	}

	if req.Email != nil && strings.TrimSpace(*req.Email) == "" {
		validationErrors = append(validationErrors, "email can not be empty")
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
