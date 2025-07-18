package dto

import (
	"github.com/google/uuid"
	"time"
)

type Response struct {
	PostID    uuid.UUID    `json:"post_id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Author    UserResponse `json:"author"`
}

type InUserResponse struct {
	PostID uuid.UUID `json:"post_id"`
	Title  string    `json:"title"`
}
