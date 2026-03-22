package dto

import (
	"event-budaya-ticketing-bcc/internal/model"
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Phone           *string    `json:"phone"`
	ProfilePhoto    *string    `json:"profile_photo"`
	Gender          *string    `json:"gender"`
	EmailVerified   bool       `json:"email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

func ToUserResponse(u *model.User) UserResponse {
	return UserResponse{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		Phone:           u.Phone,
		ProfilePhoto:    u.ProfilePhoto,
		Gender:          u.Gender,
		EmailVerified:   u.EmailVerifiedAt != nil,
		EmailVerifiedAt: u.EmailVerifiedAt,
	}
}
