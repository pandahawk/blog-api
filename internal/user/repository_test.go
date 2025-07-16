package user

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&User{}))
	sampleUsers := []User{
		{Username: "testuser1", Email: "t1@example.com"},
		{Username: "testuser2", Email: "t2@example.com"},
		{Username: "testuser3", Email: "t3@example.com"},
	}
	require.NoError(t, db.Create(&sampleUsers).Error)
	return db
}

func TestGormRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)

	users, err := repo.FindAll()

	require.NoError(t, err)
	require.Len(t, users, 3)
	assert.Equal(t, "testuser1", users[0].Username)
	assert.Equal(t, "t1@example.com", users[0].Email)

}

func TestGormRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	id := 1

	got, err := repo.FindByID(id)

	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
}

func TestGormRepository_FindByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	username := "testuser1"

	got, err := repo.FindByUsername(username)

	require.NoError(t, err)
	assert.Equal(t, username, got.Username)
}

func TestGormRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	email := "t1@example.com"

	got, err := repo.FindByEmail(email)

	require.NoError(t, err)
	assert.Equal(t, email, got.Email)
}

func TestGormRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	user := User{Username: "newuser", Email: "newemail@mail.com"}

	saved, err := repo.Create(user)

	require.NoError(t, err)
	assert.Equal(t, user.Username, saved.Username)

	saved, err = repo.Create(user)
	assert.Error(t, err)
}

func TestGormRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	user := User{Username: "newuser", Email: "newemail@mail.com"}

	updated, err := repo.Update(user)

	require.NoError(t, err)
	assert.Equal(t, user.Username, updated.Username)
}

func TestGormRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	user := User{
		ID:       1,
		Username: "testuser1",
		Email:    "t1@example.com",
	}

	err := repo.Delete(user)

	assert.NoError(t, err)
}
