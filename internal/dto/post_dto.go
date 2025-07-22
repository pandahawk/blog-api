package dto

import (
	"github.com/google/uuid"
	"time"
)

type PostResponse struct {
	PostID    uuid.UUID           `json:"post_id"`
	Title     string              `json:"title"`
	Content   string              `json:"content"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Author    UserSummaryResponse `json:"author"`
}

type PostSummaryResponse struct {
	//PostID uuid.UUID `json:"post_id"`
	//Title  string    `json:"title"`
	PostID uuid.UUID `json:"post_id" format:"uuid"`
	Title  string    `json:"title" example:"My First Post"`
}

type CreatePostRequest struct {
	Title    string    `json:"title" binding:"required"`
	Content  string    `json:"content" binding:"required"`
	AuthorID uuid.UUID `json:"author_id" binding:"required"`
}

type CreatePostInUserRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}
