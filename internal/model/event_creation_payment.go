package model

import (
	"time"

	"github.com/google/uuid"
)

type EventCreationPayment struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EventID          uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"event_id"`
	PromoterID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"promoter_id"`
	OrderID          string     `gorm:"size:100;uniqueIndex;not null" json:"order_id"`
	PaymentToken     string     `gorm:"size:500" json:"payment_token"`
	Amount           float64    `gorm:"type:numeric(12,2);not null" json:"amount"`
	Status           string     `gorm:"size:20;default:pending" json:"status"` // pending, settlement, failure, cancelled
	PaymentMethod    *string    `gorm:"size:50" json:"payment_method"`
	PaymentURL       *string    `json:"payment_url"`
	ExternalID       *string    `gorm:"size:100;index" json:"external_id"`
	ExternalResponse *string    `json:"external_response"`
	PaidAt           *time.Time `json:"paid_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`

	Event    Event `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Promoter User  `gorm:"foreignKey:PromoterID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

func (EventCreationPayment) TableName() string {
	return "event_creation_payments"
}
