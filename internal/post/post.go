package post

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"` // @swagger:strfmt uuid
	Title     string    `json:"title" gorm:"not null" binding:"required"`
	Content   string    `json:"text" gorm:"type:text;not null" binding:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	UserID    uuid.UUID `gorm:"type:char(36);not null"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
