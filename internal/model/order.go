package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	EventID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"event_id"`
	Quantity   int        `gorm:"not null" json:"quantity"`
	UnitPrice  float64    `gorm:"type:numeric(12,2);not null" json:"unit_price"`
	ServiceFee float64    `gorm:"type:numeric(12,2);default:0" json:"service_fee"`
	TotalPrice float64    `gorm:"type:numeric(12,2);default:0" json:"total_price"`
	Status     string     `gorm:"size:20;default:pending" json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	User       User       `gorm:"foreignKey:UserID" json:"-"`
	Event      Event      `gorm:"foreignKey:EventID" json:"-"`
}

func (Order) TableName() string {
	return "orders"
}
