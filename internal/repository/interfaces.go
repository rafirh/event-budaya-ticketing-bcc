package repository

import (
	"time"

	"event-budaya-ticketing-bcc/internal/model"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByID(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
}

type PersonalAccessTokenRepository interface {
	Create(token *model.PersonalAccessToken) error
	FindByToken(token string) (*model.PersonalAccessToken, error)
	DeleteByToken(token string) error
	DeleteByUserID(userID uuid.UUID) error
	UpdateLastUsed(id uint, t time.Time) error
}

type EventCategoryRepository interface {
	FindAll(limit, offset int, sortBy, sortOrder string) ([]model.EventCategory, int64, error)
}

type EventRepository interface {
	FindAll(search, categoryID, sortBy, sortOrder string, limit, offset int) ([]model.Event, int64, error)
	FindBySlug(slug string) (*model.Event, error)
}
