package model

import (
	"time"

	"github.com/google/uuid"
)

type PersonalAccessToken struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Token      string     `gorm:"not null;uniqueIndex" json:"-"`
	Name       string     `gorm:"size:100;default:'default'" json:"name"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (PersonalAccessToken) TableName() string {
	return "personal_access_tokens"
}
