package user

import (
	"log"
	"sync"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=user
type Repository interface {
	GetAllUsers() ([]User, error)
	GetUserById(id string) (User, bool)
}

type repository struct {
	users map[string]User
	mu    sync.RWMutex
}

func (r *repository) GetUserById(id string) (User, bool) {
	user, ok := r.users[id]
	return user, ok
}

func (r *repository) GetAllUsers() ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]User, 0)
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *repository) InitSampleData() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = map[string]User{
		"1": {ID: "1", Username: "alice", Email: "alice@example.com"},
		"2": {ID: "2", Username: "bob", Email: "bob@example.com"},
		"3": {ID: "3", Username: "charlie", Email: "charlie@example.com"},
	}
	log.Println("Sample Data initialized")

}

func NewRepository() Repository {
	return &repository{users: make(map[string]User)}
}

func NewDevRepository() Repository {
	repo := &repository{
		users: make(map[string]User),
	}
	repo.InitSampleData()
	return repo
}
