package user

import (
	"github.com/pandahawk/blog-api/internal/database"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"github.com/pandahawk/blog-api/internal/shared/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepository_FindAll(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewRepository(db)

	users, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Equal(t, len(testdata.SampleUsers), len(users))
	for i := range users {
		assert.Equal(t, testdata.SampleUsers[i].ID, users[i].ID)
		assert.Equal(t, testdata.SampleUsers[i].Username, users[i].Username)
		assert.Equal(t, testdata.SampleUsers[i].Email, users[i].Email)
	}
}

func TestRepository_FindByID(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewRepository(db)
	got, err := repo.FindByID(testdata.Alice.ID)

	assert.NoError(t, err)
	assert.Equal(t, testdata.Alice.ID, got.ID)
	assert.Equal(t, testdata.Alice.Username, got.Username)
}

func TestRepository_FindByUsername(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewRepository(db)
	username := "bob"

	got, err := repo.FindByUsername(username)

	assert.NoError(t, err)
	assert.Equal(t, testdata.Bob.Username, got.Username)
}

func TestRepository_FindByEmail(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewRepository(db)
	got, err := repo.FindByEmail(testdata.Caren.Email)

	require.NoError(t, err)
	assert.Equal(t, testdata.Caren.Email, got.Email)
}

func TestRepository_Create(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewRepository(db)
	user := model.NewUser("newuser", "newemail@mail.com")

	got, err := repo.Create(user)

	assert.NoError(t, err)
	assert.Equal(t, user.Username, got.Username)
	assert.Equal(t, user.Email, got.Email)

	got, err = repo.Create(user)
	assert.Error(t, err)
}

func TestRepository_Update(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewRepository(db)
	aliceUpdate := *testdata.Alice
	aliceUpdate.Username = "newname"
	aliceUpdate.Email = "newmail@example.com"

	got, err := repo.Update(&aliceUpdate)

	require.NoError(t, err)
	assert.Equal(t, aliceUpdate.Username, got.Username)
}

func TestRepository_Delete(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewRepository(db)
	err := repo.Delete(testdata.Alice)
	assert.NoError(t, err)
}
