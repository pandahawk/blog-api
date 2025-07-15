package user

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}
