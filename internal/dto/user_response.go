package dto

import (
	"event-budaya-ticketing-bcc/internal/model"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        *string   `json:"phone"`
	Role         string    `json:"role"`
	ProfilePhoto *string   `json:"profile_photo"`
	Gender       *string   `json:"gender"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

func ToUserResponse(u *model.User) UserResponse {
	return UserResponse{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Phone:        u.Phone,
		Role:         u.Role,
		ProfilePhoto: u.ProfilePhoto,
		Gender:       u.Gender,
	}
}
