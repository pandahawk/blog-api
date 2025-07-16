package user

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
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

// todo: fix bug: this test fails when all are run. leakage?
func TestGormRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormRepository(db)
	id := 1
	got, err := repo.FindByID(id)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
}
