package post

import (
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

var sampleUsers []*model.User
var samplePosts []*model.Post

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Post{}))
	sampleUsers = []*model.User{
		{ID: uuid.MustParse("4e7b1a2c-9d3e-4f5a-8b6c-7d9e0f1a2b3c"),
			Username: "testuser1",
			Email:    "t1@example.com"},
		{ID: uuid.MustParse("a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"),
			Username: "testuser2", Email: "t2@example.com"},
	}

	samplePosts = []*model.Post{
		{
			ID:        uuid.MustParse("7f8e9d0c-1b2a-4345-6789-abcdef012345"),
			Title:     "First Post",
			Content:   "This is the content of the first post.",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			UserID:    sampleUsers[0].ID,
			User:      sampleUsers[0],
		},
		{
			ID:        uuid.MustParse("c8d9e0f1-2a3b-4c4d-5e6f-7a8b9c0d1e2f"),
			Title:     "Second Post",
			Content:   "This is the content of the second post.",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			UserID:    sampleUsers[0].ID,
			User:      sampleUsers[0],
		},
		{
			ID:        uuid.MustParse("2b3c4d5e-6f7a-4b8c-9d0e-1f2a3b4c5d6e"),
			Title:     "Second Post",
			Content:   "This is the content of the second post.",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			UserID:    sampleUsers[1].ID,
			User:      sampleUsers[1],
		},
	}

	require.NoError(t, db.Create(sampleUsers).Error)
	require.NoError(t, db.Create(samplePosts).Error)
	return db
}

func TestRepository_FindAll(t *testing.T) {
	t.Run("should find all posts", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)

		posts, err := repo.FindAll()
		require.NoError(t, err)
		require.Len(t, posts, 3)
		assert.Equal(t, samplePosts[0].ID, posts[0].ID)

	})
}

func TestRepository_FindByID(t *testing.T) {
	t.Run("should find post by ID", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewRepository(db)
		id := samplePosts[1].ID

		post, err := repo.FindByID(id)

		require.NoError(t, err)
		assert.Equal(t, samplePosts[1], post)
	})
}
