package user

import (
	"log"
	"sort"
	"sync"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=user
type Repository interface {
	FindAll() ([]User, error)
	FindByID(id int) (User, bool)
	Save(user User) (User, error)
}

type repository struct {
	users     map[int]User
	mu        sync.RWMutex
	idCounter int
}

func (r *repository) FindByID(id int) (User, bool) {
	user, ok := r.users[id]
	return user, ok
}

func (r *repository) FindAll() ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]User, 0)
	for _, user := range r.users {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	return users, nil
}

func (r *repository) InitSampleData() {
	r.users = map[int]User{
		101: {ID: 101, Username: "alice", Email: "alice@example.com"},
		102: {ID: 102, Username: "bob", Email: "bob@example.com"},
		103: {ID: 103, Username: "charlie", Email: "charlie@example.com"},
	}
	r.idCounter = 100 + len(r.users)
	log.Println("Sample Data initialized")
}

func (r *repository) Save(u User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.idCounter++
	u.ID = r.idCounter
	r.users[u.ID] = u
	return u, nil
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
