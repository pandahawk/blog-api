package mapper

import (
	"github.com/pandahawk/blog-api/internal/dto"
	"github.com/pandahawk/blog-api/internal/post"
	"github.com/pandahawk/blog-api/internal/user"
)

func fromPostAndUser(p post.Post, author user.User) dto.PostResponse {
	return dto.PostResponse{
		PostID:    p.ID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Author:    FromUserInPostResponse(author),
	}
}

func fromPostInUserResponse(p post.Post) dto.PostSummaryResponse {
	return dto.PostSummaryResponse{
		PostID: p.ID,
		Title:  p.Title,
	}
}
