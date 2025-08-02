package post

import (
	"github.com/google/uuid"
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

	posts, err := repo.FindAll()
	require.NoError(t, err)
	require.Len(t, posts, len(testdata.SamplePosts))

	for i := range posts {
		assert.Equal(t, testdata.SamplePosts[i].ID, posts[i].ID)
		assert.Equal(t, testdata.SamplePosts[i].Title, posts[i].Title)
		assert.Equal(t, testdata.SamplePosts[i].Content, posts[i].Content)
		assert.Equal(t, testdata.SamplePosts[i].UserID, posts[i].UserID)
	}
}

func TestRepository_FindByID(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewRepository(db)

	post, err := repo.FindByID(testdata.PostIDs[0])

	require.NoError(t, err)
	assert.Equal(t, testdata.Post1.ID, post.ID)
	assert.Equal(t, testdata.Post1.UserID, post.UserID)
	assert.Equal(t, testdata.Post1.Title, post.Title)
	assert.Equal(t, testdata.Post1.Content, post.Content)
}

func TestRepository_Create(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewRepository(db)
	post := model.NewPost("a new got",
		"content of a new got",
		testdata.Caren.ID)

	got, err := repo.Create(post)

	require.NoError(t, err)
	assert.Equal(t, "content of a new got", got.Content)
	assert.Equal(t, "a new got", got.Title)
	assert.Equal(t, testdata.Caren.ID, got.UserID)

	got, err = repo.Create(post)
	assert.Error(t, err)
}

func TestRepository_CreateWithouUser(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewRepository(db)

	_, err := repo.Create(model.NewPost("test", "test", uuid.Nil))

	assert.Error(t, err)

}

func TestRepository_Delete(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewRepository(db)

	err := repo.Delete(testdata.Post1)
	require.NoError(t, err)
}

func TestRepository_Update(t *testing.T) {
	db := database.SetupTestDB(t)
	repo := NewRepository(db)
	updatedPost := *testdata.Post1 // This is a value copy (shallow)
	updatedPost.Title = "new title"
	updatedPost.Content = "new content"

	post, err := repo.Update(&updatedPost)

	require.NoError(t, err)
	assert.Equal(t, updatedPost.ID, post.ID)
	assert.Equal(t, updatedPost.Title, "new title")
	assert.Equal(t, updatedPost.Content, "new content")

}
