package user

import (
	"fmt"
	"log"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=user

type Service interface {
	GetUser(id int) (User, error)
	CreateUser(user User) (User, error)
	GetAllUsers() []User
}

type service struct {
	repo Repository
}

func (s *service) CreateUser(user User) (User, error) {
	return user, nil
}

func (s *service) GetUser(id int) (User, error) {
	user, ok := s.repo.GetUserByID(id)
	if !ok {
		return User{}, fmt.Errorf("user with ID %v not found", id)
	}
	return user, nil
}

func (s *service) GetAllUsers() []User {

	users, err := s.repo.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}
	return users
}

func NewService(r Repository) Service {
	return &service{repo: r}
}
