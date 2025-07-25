package user

import (
	"github.com/google/uuid"
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email" binding:"omitempty,email"`
}

type Response struct {
	UserID   uuid.UUID              `json:"user_id" swaggertype:"string"`
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	JoinedAt time.Time              `json:"joined_at"`
	Posts    []*PostSummaryResponse `json:"posts"`
}

type PostSummaryResponse struct {
	PostID uuid.UUID `json:"post_id" format:"uuid"`
	Title  string    `json:"title" example:"My First Post"`
}
