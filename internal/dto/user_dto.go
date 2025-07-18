package dto

import (
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email" binding:"email"`
}

type UserResponse struct {
	ID       uuid.UUID        `json:"author_id"`
	Username string           `json:"username"`
	Email    string           `json:"email"`
	Posts    []InUserResponse `json:"posts"`
}

type InPostResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

//todo: think of author as full response with []posts and user as just user
