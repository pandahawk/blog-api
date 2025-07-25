package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=user

type Service interface {
	GetUser(id uuid.UUID) (*model.User, error)
	CreateUser(req *CreateUserRequest) (*model.User, error)
	GetUsers() ([]*model.User, error)
	UpdateUser(id uuid.UUID, req *UpdateUserRequest) (*model.User, error)
	DeleteUser(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func validateUsernameFormat(username string) error {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]{3,}$`, username)

	if !matched {
		return apperrors.NewInvalidInputError("invalid username: must be" +
			" alphanumeric, at least 3 character")
	}

	if _, err := strconv.ParseFloat(username, 64); err == nil {
		return apperrors.NewInvalidInputError(
			"invalid username: must not be a number")
	}

	letters := 0
	for _, r := range username {
		if unicode.IsLetter(r) {
			letters++
		}
	}

	if letters < 2 {
		return apperrors.NewInvalidInputError(
			"invalid username: must have at least two letters")
	}

	return nil
}

func (s *service) CreateUser(req *CreateUserRequest) (*model.User, error) {

	if err := validateUsernameFormat(req.Username); err != nil {
		return nil, err
	}

	user, err := s.repo.Create(model.NewUser(req.Username, req.Email))
	if err != nil {
		if strings.Contains(err.Error(),
			`violates unique constraint "uni_users_username"`) {
			return nil, apperrors.NewDuplicateError("username")
		}
		if strings.Contains(err.Error(),
			`violates unique constraint "uni_users_email"`) {
			return nil, apperrors.NewDuplicateError("email")
		}
		return nil, err
	}
	return user, nil
}

func (s *service) UpdateUser(id uuid.UUID, req *UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("user", id)
	}

	if req.Username != nil {
		if err := validateUsernameFormat(*req.Username); err != nil {
			return nil, err
		}
		if _, err := s.repo.FindByUsername(*req.Username); err == nil {
			return nil, apperrors.NewDuplicateError("username already exists")
		}

		user.Username = *req.Username
	}

	if req.Email != nil {
		if _, err := s.repo.FindByEmail(*req.Email); err == nil {
			return nil, apperrors.NewDuplicateError("email")
		}
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

func (s *service) GetUser(id uuid.UUID) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("user", id)
	}
	return user, nil
}

func (s *service) GetUsers() ([]*model.User, error) {

	users, err := s.repo.FindAll()
	if err != nil {
		return nil, errors.New("failed to get all users")
	}
	return users, nil
}

func NewService(r Repository) Service {
	return &service{repo: r}
}
