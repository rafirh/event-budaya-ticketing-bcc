package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name         string     `gorm:"size:150;not null" json:"name"`
	Email        string     `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Password     string     `gorm:"not null" json:"-"`
	Phone        *string    `gorm:"size:20" json:"phone"`
	Role         string     `gorm:"size:20;not null" json:"role"`
	ProfilePhoto *string    `json:"profile_photo"`
	Gender       *string    `gorm:"size:10" json:"gender"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
