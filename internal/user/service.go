package user

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=user

type Service interface {
	GetUser(id int) (User, error)
	CreateUser(req CreateUserRequest) (User, error)
	GetAllUsers() []User
	UpdateUser(id int, req UpdateUserRequest) (User, error)
	DeleteUser(id int) error
}

type service struct {
	repo Repository
}

func (s *service) CreateUser(req CreateUserRequest) (User, error) {
	_, found := s.repo.FindByUsername(req.Username)
	if found {
		return User{}, errors.New("username already exists")
	}
	_, found = s.repo.FindByEmail(req.Email)
	if found {
		return User{}, errors.New("email already exists")
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
		return User{}, fmt.Errorf("user %d not found", id)
	}

	// Only update provided fields
	if req.Username != nil && strings.TrimSpace(*req.Username) == "" {
		return User{}, errors.New("username cannot be empty")
	}

	if req.Email != nil && strings.TrimSpace(*req.Email) == "" {
		return User{}, errors.New("email cannot be empty")
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
	//TODO implement me
	panic("implement me")
}

func (s *service) GetUser(id int) (User, error) {
	user, ok := s.repo.FindByID(id)
	if !ok {
		return User{}, fmt.Errorf("user with ID %v not found", id)
	}
	return user, nil
}

func (s *service) GetAllUsers() []User {

	users, err := s.repo.FindAll()
	if err != nil {
		log.Fatal(err)
	}
	return users
}

func NewService(r Repository) Service {
	return &service{repo: r}
}
