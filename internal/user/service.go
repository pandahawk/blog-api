package user

import (
	"errors"
	"fmt"
	"log"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=user

type Service interface {
	GetUser(id int) (User, error)
	CreateUser(user User) (User, error)
	GetAllUsers() []User
	UpdateUser(id int, user User) (User, error)
	DeleteUser(id int) error
}

type service struct {
	repo Repository
}

func (s *service) DeleteUser(id int) error {
	//TODO implement me
	panic("implement me")
}

func (s *service) UpdateUser(id int, user User) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) CreateUser(user User) (User, error) {

	if user.Username == "" || user.Email == "" {
		return User{}, errors.New("username or email is empty")
	}

	users, err := s.repo.FindAll()
	if err != nil {
		return User{}, err
	}

	for _, u := range users {
		if u.Username == user.Username {
			return User{}, errors.New("username already exists")
		}
		if u.Email == user.Email {
			return User{}, errors.New("email already exists")
		}
	}
	return s.repo.Save(user)
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
