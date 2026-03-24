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
	DeleteByID(id string) error
}

type EmailVerificationTokenRepository interface {
	Create(token *model.EmailVerificationToken) error
	FindByTokenHash(tokenHash string) (*model.EmailVerificationToken, error)
	FindLatestByUserID(userID uuid.UUID) (*model.EmailVerificationToken, error)
	DeleteByUserID(userID uuid.UUID) error
	MarkUsed(id uint, usedAt time.Time) error
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
	Create(event *model.Event) error
	FindAll(search, categoryID, sortBy, sortOrder string, limit, offset int) ([]model.Event, int64, error)
	FindByID(id string) (*model.Event, error)
	FindBySlug(slug string) (*model.Event, error)
	FindByPromoterID(promoterID uuid.UUID, limit, offset int) ([]model.Event, int64, error)
	Update(event *model.Event) error
}

type EventCommentRepository interface {
	Create(comment *model.EventComment) error
	FindByID(id string) (*model.EventComment, error)
	FindByEventID(eventID string) ([]model.EventComment, error)
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

type PromoterWalletRepository interface {
	Create(wallet *model.PromotorWallet) error
	FindByPromoterID(promoterID string) (*model.PromotorWallet, error)
	Update(wallet *model.PromotorWallet) error
}

type WalletTransactionRepository interface {
	Create(transaction *model.WalletTransaction) error
	FindByWalletID(walletID string) ([]model.WalletTransaction, error)
}

type FeeRepository interface {
	FindByType(feeType string) (*model.Fee, error)
	FindAll() ([]model.Fee, error)
}

type AdminWalletRepository interface {
	FindOrCreate() (*model.AdminWallet, error)
	Update(wallet *model.AdminWallet) error
	FindByID(id string) (*model.AdminWallet, error)
}

type PromoterTransactionHistoryRepository interface {
	Create(transaction *model.PromoterTransactionHistory) error
	FindByPromoterID(promoterID uuid.UUID, limit, offset int) ([]model.PromoterTransactionHistory, int64, error)
}

type EventCreationPaymentRepository interface {
	Create(payment *model.EventCreationPayment) error
	FindByOrderID(orderID string) (*model.EventCreationPayment, error)
	FindByEventID(eventID uuid.UUID) (*model.EventCreationPayment, error)
	Update(payment *model.EventCreationPayment) error
}
