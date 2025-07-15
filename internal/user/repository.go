package user

import (
	"log"
	"sort"
	"sync"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=user
type Repository interface {
	FindAll() ([]User, bool)
	FindByID(id int) (User, bool)
	FindByUsername(username string) (User, bool)
	FindByEmail(email string) (User, bool)
	Save(user User) (User, bool)
	Update(user User) (User, bool)
	Delete(user User) bool
}

type repository struct {
	users     map[int]User
	mu        sync.RWMutex
	idCounter int
}

func (r *repository) Delete(user User) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, user.ID)
	return true
}

func (r *repository) Update(user User) (User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return user, true
}

func (r *repository) FindByUsername(username string) (User, bool) {
	for _, user := range r.users {
		if user.Username == username {
			return user, true
		}
	}
	return User{}, false
}

func (r *repository) FindByEmail(email string) (User, bool) {
	for _, user := range r.users {
		if user.Email == email {
			return user, true
		}
	}
	return User{}, false
}

func (r *repository) FindByID(id int) (User, bool) {
	user, ok := r.users[id]
	return user, ok
}

func (r *repository) FindAll() ([]User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]User, 0)
	for _, user := range r.users {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	return users, true
}

func (r *repository) InitSampleData() {
	r.users = map[int]User{
		1001: {ID: 1001, Username: "alice", Email: "alice@example.com"},
		1002: {ID: 1002, Username: "bob", Email: "bob@example.com"},
		1003: {ID: 1003, Username: "charlie", Email: "charlie@example.com"},
		1004: {ID: 1004, Username: "diana", Email: "diana@example.com"},
		1005: {ID: 1005, Username: "eve", Email: "eve@example.com"},
		1006: {ID: 1006, Username: "frank", Email: "frank@example.com"},
		1007: {ID: 1007, Username: "grace", Email: "grace@example.com"},
		1008: {ID: 1008, Username: "heidi", Email: "heidi@example.com"},
		1009: {ID: 1009, Username: "ivan", Email: "ivan@example.com"},
		1010: {ID: 1010, Username: "judy", Email: "judy@example.com"},
		1011: {ID: 1011, Username: "kevin", Email: "kevin@example.com"},
		1012: {ID: 1012, Username: "lisa", Email: "lisa@example.com"},
		1013: {ID: 1013, Username: "mike", Email: "mike@example.com"},
		1014: {ID: 1014, Username: "nancy", Email: "nancy@example.com"},
		1015: {ID: 1015, Username: "oliver", Email: "oliver@example.com"},
		1016: {ID: 1016, Username: "patricia", Email: "patricia@example.com"},
		1017: {ID: 1017, Username: "quinn", Email: "quinn@example.com"},
		1018: {ID: 1018, Username: "randy", Email: "randy@example.com"},
		1019: {ID: 1019, Username: "sara", Email: "sara@example.com"},
		1020: {ID: 1020, Username: "tom", Email: "tom@example.com"},
	}
	r.idCounter = 1000 + len(r.users)
	log.Println("Sample Data initialized")
}

func (r *repository) Save(u User) (User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.idCounter++
	u.ID = r.idCounter
	r.users[u.ID] = u
	return u, true
}

func NewRepository() Repository {
	return &repository{
		users:     make(map[int]User),
		idCounter: 1000}
}

func NewDevRepository() Repository {
	repo := &repository{
		users:     make(map[int]User),
		idCounter: 1000,
	}
	repo.InitSampleData()
	return repo
}
