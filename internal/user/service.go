package user

import "fmt"

type UserService interface {
	GetAllUsers() string
	GetUser(id string) string
	CreateUser() string
	UpdateUser(id string) string
	DeleteUser(id string) string
}

type simpleService struct {
}

func (s *simpleService) GetUser(id string) string {
	return fmt.Sprintf("get user %s", id)
}

func (s *simpleService) CreateUser() string {
	return "create new user"
}

func (s *simpleService) UpdateUser(id string) string {
	return fmt.Sprintf("update user %s", id)
}

func (s *simpleService) DeleteUser(id string) string {
	return fmt.Sprintf("delete user %s", id)
}

func NewSimpleService() UserService {
	return &simpleService{}
}

func (s *simpleService) GetAllUsers() string {
	return "get all users"
}
