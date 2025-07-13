package user

import (
	"log"
	"sync"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=user
type Repository interface {
	GetAllUsers() ([]User, error)
	GetUserByID(id int) (User, bool)
	SaveUser(user User) (User, error)
}

type repository struct {
	users     map[int]User
	mu        sync.RWMutex
	idCounter int
}

func (r *repository) SaveUser(user User) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetUserByID(id int) (User, bool) {
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

/*
todo: write a save method and call this instead of manually reassigning
idCounter
*/
func (r *repository) InitSampleData() {
	r.Save(User{Username: "alice", Email: "alice@example.com"})
	r.Save(User{Username: "bob", Email: "bob@example.com"})
	r.Save(User{Username: "charlie", Email: "charlie@example.com"})
	log.Println("Sample Data initialized")
}

func (r *repository) Save(u User) User {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.idCounter++
	u.ID = r.idCounter
	r.users[u.ID] = u
	return u
}

func NewRepository() Repository {
	return &repository{
		users:     make(map[int]User),
		idCounter: 100}
}

func NewDevRepository() Repository {
	repo := &repository{
		users:     make(map[int]User),
		idCounter: 100,
	}
	repo.InitSampleData()
	return repo
}
