package user

import (
	"errors"
	"log"
	"sort"
	"sync"
)

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=user
type Repository interface {
	FindAll() ([]User, error)
	FindByID(id int) (User, error)
	FindByUsername(username string) (User, error)
	FindByEmail(email string) (User, error)
	Save(user User) (User, error)
	Update(user User) (User, error)
	Delete(user User) error
}

type repository struct {
	users     map[int]User
	mu        sync.RWMutex
	idCounter int
}

func (r *repository) Delete(user User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, user.ID)
	return nil
}

func (r *repository) Update(user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return user, nil
}

func (r *repository) FindByUsername(username string) (User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}

func (r *repository) FindByEmail(email string) (User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}

func (r *repository) FindByID(id int) (User, error) {
	user, ok := r.users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}
	return user, nil
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
		1001: {ID: 1001, Username: "blogger01", Email: "blogger01@example.com"},
		1002: {ID: 1002, Username: "blogger02", Email: "blogger02@example.com"},
		1003: {ID: 1003, Username: "blogger03", Email: "blogger03@example.com"},
		1004: {ID: 1004, Username: "blogger04", Email: "blogger04@example.com"},
		1005: {ID: 1005, Username: "blogger05", Email: "blogger05@example.com"},
		1006: {ID: 1006, Username: "blogger06", Email: "blogger06@example.com"},
		1007: {ID: 1007, Username: "blogger07", Email: "blogger07@example.com"},
		1008: {ID: 1008, Username: "blogger08", Email: "blogger08@example.com"},
		1009: {ID: 1009, Username: "blogger09", Email: "blogger09@example.com"},
		1010: {ID: 1010, Username: "blogger10", Email: "blogger10@example.com"},
		1011: {ID: 1011, Username: "blogger11", Email: "blogger11@example.com"},
		1012: {ID: 1012, Username: "blogger12", Email: "blogger12@example.com"},
		1013: {ID: 1013, Username: "blogger13", Email: "blogger13@example.com"},
		1014: {ID: 1014, Username: "blogger14", Email: "blogger14@example.com"},
		1015: {ID: 1015, Username: "blogger15", Email: "blogger15@example.com"},
		1016: {ID: 1016, Username: "blogger16", Email: "blogger16@example.com"},
		1017: {ID: 1017, Username: "blogger17", Email: "blogger17@example.com"},
		1018: {ID: 1018, Username: "blogger18", Email: "blogger18@example.com"},
		1019: {ID: 1019, Username: "blogger19", Email: "blogger19@example.com"},
		1020: {ID: 1020, Username: "blogger20", Email: "blogger20@example.com"},
	}
	r.idCounter = 1000 + len(r.users)
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

//todo: replace dev repository with this one when postgres is connected
//func NewRepository() Repository {
//	return &repository{
//		users:     make(map[int]User),
//		idCounter: 1000}
//}

func NewDevRepository() Repository {
	repo := &repository{
		users:     make(map[int]User),
		idCounter: 1000,
	}
	repo.InitSampleData()
	return repo
}
