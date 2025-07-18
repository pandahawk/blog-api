package post

import (
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/user"
	"time"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"` // @swagger:strfmt uuid
	Title     string    `json:"title" gorm:"not null" binding:"required"`
	Text      string    `json:"text" gorm:"not null" binding:"required"`
	Author    user.User `json:"author" gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}
