package model

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderID        uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"order_id"`
	PaymentMethod  *string    `gorm:"size:50" json:"payment_method"`
	PaymentGateway *string    `gorm:"size:50" json:"payment_gateway"`
	Amount         float64    `gorm:"type:numeric(12,2)" json:"amount"`
	Status         string     `gorm:"size:20" json:"status"`
	PaidAt         *time.Time `json:"paid_at"`
	CreatedAt      time.Time  `json:"created_at"`
	Order          Order      `gorm:"foreignKey:OrderID" json:"-"`
}

func (Payment) TableName() string {
	return "payments"
}
