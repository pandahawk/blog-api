package mapper

import (
	"github.com/pandahawk/blog-api/internal/dto"
	"github.com/pandahawk/blog-api/internal/user"
)

func FromUser(u *user.User) dto.UserResponse {

	posts := make([]dto.PostSummaryResponse, len(u.Posts))
	for i, p := range u.Posts {
		posts[i] = fromPostInUserResponse(p)
	}
	return dto.UserResponse{
		UserID:   u.ID,
		Username: u.Username,
		Email:    u.Email,
		Posts:    posts,
		JoinedAt: u.CreatedAt,
	}
}

func FromUserInPostResponse(u *user.User) dto.UserSummaryResponse {
	return dto.UserSummaryResponse{
		UserID:   u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}
