package repository

import (
	"event-budaya-ticketing-bcc/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
}
