package domain

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

type UserResponse struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Phone        *string    `json:"phone"`
	Role         string     `json:"role"`
	ProfilePhoto *string    `json:"profile_photo"`
	Gender       *string    `json:"gender"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Phone:        u.Phone,
		Role:         u.Role,
		ProfilePhoto: u.ProfilePhoto,
		Gender:       u.Gender,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

type RegisterRequest struct {
	Name     string  `json:"name" validate:"required,min=2,max=150"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Phone    *string `json:"phone" validate:"omitempty,max=20"`
	Role     string  `json:"role" validate:"omitempty,oneof=user promotor"`
	Gender   *string `json:"gender" validate:"omitempty,oneof=male female other"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
