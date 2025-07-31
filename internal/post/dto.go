package post

import (
	"github.com/google/uuid"
	"time"
)

type CreatePostRequest struct {
	Title    string    `json:"title" binding:"required"`
	Content  string    `json:"content" binding:"required"`
	AuthorID uuid.UUID `json:"author_id" binding:"required"`
}

type UpdatePostRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type Response struct {
	PostID    uuid.UUID           `json:"post_id"`
	Title     string              `json:"title"`
	Content   string              `json:"content"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Author    UserSummaryResponse `json:"author"`
}

type UserSummaryResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}
