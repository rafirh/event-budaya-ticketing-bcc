package repository

import (
	"time"

	"event-budaya-ticketing-bcc/internal/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
}

type PersonalAccessTokenRepository interface {
	Create(token *domain.PersonalAccessToken) error
	FindByToken(token string) (*domain.PersonalAccessToken, error)
	DeleteByToken(token string) error
	DeleteByUserID(userID uuid.UUID) error
	UpdateLastUsed(id uint, t time.Time) error
}
