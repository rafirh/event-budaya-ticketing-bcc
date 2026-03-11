package repository

import (
	"event-budaya-ticketing-bcc/internal/domain"
)

// UserRepository interface defines user repository methods
type UserRepository interface {
	Create(user *domain.User) error
	FindAll() ([]domain.User, error)
	FindByID(id uint) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
}
