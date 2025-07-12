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

func (s *simpleService) GetAllUsers() string {
	return "Get All users"
}

func NewSimpleService() UserService {
	return &simpleService{}
}

func (s *simpleService) GetUser(id string) string {
	return fmt.Sprintf("Get user %s", id)
}

func (s *simpleService) CreateUser() string {
	return "Create new user"
}

func (s *simpleService) UpdateUser(id string) string {
	return fmt.Sprintf("Update user %s", id)
}

func (s *simpleService) DeleteUser(id string) string {
	return fmt.Sprintf("Delete user %s", id)
}
