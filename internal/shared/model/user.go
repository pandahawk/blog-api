package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" example:"b9e69a63-4f4b-4ea7-8c71-3b73fe62e6d7"`
	Username  string    `json:"username" gorm:"unique;not null" example:"mike"`
	Email     string    `json:"email" gorm:"unique;not null" example:"mike@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2025-07-18T15:04:05Z"`
	Posts     []*Post   `gorm:"foreignKey:UserID"`
}

func NewUser(username string, email string) *User {
	return &User{Username: username, Email: email}
}

//goland:noinspection GoUnusedParameter
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
