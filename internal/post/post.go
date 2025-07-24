package post

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" example:"4e76b320-d5b7-4a0a-bb0f-2049fe6a91a7"`
	Title     string    `json:"title" gorm:"not null" binding:"required" example:"My First Post"`
	Content   string    `json:"content" gorm:"type:text;not null" binding:"required" example:"My First Post Content"`
	CreatedAt time.Time `json:"created_at" gorm:"not null" example:"2025-07-18T15:04:05Z"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null" example:"2025-08-19T15:04:05Z"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;constraint:OnDelete:CASCADE;" example:"5e76b320-d5b7-4a0a-bb0f-2049fe3a91a4"`
}

func NewPost(title string, content string, authorID uuid.UUID) *Post {
	return &Post{
		Title:   title,
		Content: content,
		UserID:  authorID,
	}
}

//goland:noinspection GoExportedElementShouldHaveComment
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
