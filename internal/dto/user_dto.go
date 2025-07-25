package dto

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

type UserResponse struct {
	UserID   uuid.UUID             `json:"user_id" swaggertype:"string"`
	Username string                `json:"username"`
	Email    string                `json:"email"`
	JoinedAt time.Time             `json:"joined_at"`
	Posts    []PostSummaryResponse `json:"posts"`
}

type UserSummaryResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}
