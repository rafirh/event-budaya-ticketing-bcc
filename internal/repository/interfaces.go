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

// ProductRepository interface defines product repository methods
type ProductRepository interface {
	Create(product *domain.Product) error
	FindAll() ([]domain.Product, error)
	FindByID(id uint) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uint) error
}
