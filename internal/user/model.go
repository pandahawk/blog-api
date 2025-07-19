package user

import (
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/post"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID   `gorm:"type:char(36);primaryKey"` // @swagger:strfmt uuid
	Username  string      `json:"username" gorm:"unique;not null"`
	Email     string      `json:"email" gorm:"unique;not null"`
	CreatedAt time.Time   `json:"created_at"`
	Posts     []post.Post `gorm:"foreignKey:UserID"`
}

//goland:noinspection GoUnusedParameter
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
