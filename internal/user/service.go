package user

import (
	"fmt"
	"log"
)

type Service interface {
	//GetAllUsers() string
	GetUser(id string) (User, error)
	//CreateUser() string
	//UpdateUser(id string) string
	//DeleteUser(id string) string
	GetAllUsers() []User
}

//type simpleService struct {
//}
//
//func (s *simpleService) GetAllUsers() string {
//	return "Get All users"
//}
//
//func NewSimpleService() Service {
//	return &simpleService{}
//}
//
//func (s *simpleService) GetUser(id string) string {
//	return fmt.Sprintf("Get user %s", id)
//}
//
//func (s *simpleService) CreateUser() string {
//	return "Create new user"
//}
//
//func (s *simpleService) UpdateUser(id string) string {
//	return fmt.Sprintf("Update user %s", id)
//}
//
//func (s *simpleService) DeleteUser(id string) string {
//	return fmt.Sprintf("Delete user %s", id)
//}

//go:generate mockgen -source=service.go -destination=service_mock.go -package=user
type service struct {
	repo Repository
}

func (s *service) GetUser(id string) (User, error) {
	user, ok := s.repo.GetUserById(id)
	if !ok {
		return User{}, fmt.Errorf("user with ID %s not found", id)
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
