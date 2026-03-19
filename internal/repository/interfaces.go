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
	FindByID(id string) (*model.Event, error)
	FindBySlug(slug string) (*model.Event, error)
	Update(event *model.Event) error
}

type OrderRepository interface {
	Create(order *model.Order) error
	FindByID(id string) (*model.Order, error)
	FindByIDWithRelations(id string) (*model.Order, error)
	FindByUserID(userID string) ([]model.Order, error)
	Update(order *model.Order) error
}

type TicketRepository interface {
	CreateBatch(tickets []model.Ticket) error
	FindByOrderID(orderID string) ([]model.Ticket, error)
	FindByID(id string) (*model.Ticket, error)
	FindByUserID(userID string) ([]model.Ticket, error)
}

type PaymentRepository interface {
	Create(payment *model.Payment) error
	FindByOrderID(orderID string) (*model.Payment, error)
	Update(payment *model.Payment) error
}
