package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required" example:"mike"`
	Email    string `json:"email" binding:"required,email" example:"mike@example.com"`
}

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email" binding:"omitempty,email"`
}

type UserResponse struct {
	UserID   uuid.UUID             `json:"user_id" example:"b9e69a63-4f4b-4ea7-8c71-3b73fe62e6d7" swaggertype:"string" format:"uuid"`
	Username string                `json:"username" example:"mike"`
	Email    string                `json:"email" example:"mike@example.com"`
	JoinedAt time.Time             `json:"joined_at" example:"2025-07-18T15:04:05Z"`
	Posts    []PostSummaryResponse `json:"posts"`
}

type UserSummaryResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

//todo: think of author as full response with []posts and user as just user
