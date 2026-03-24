package model

import (
	"time"

	"github.com/google/uuid"
)

type PromoterTransactionHistory struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PromoterID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"promoter_id"`
	TransactionType string     `gorm:"size:50;not null;index" json:"transaction_type"`
	Direction       string     `gorm:"size:10;not null" json:"direction"`
	Amount          float64    `gorm:"type:numeric(12,2);not null" json:"amount"`
	BalanceBefore   *float64   `gorm:"type:numeric(14,2)" json:"balance_before"`
	BalanceAfter    *float64   `gorm:"type:numeric(14,2)" json:"balance_after"`
	ReferenceType   *string    `gorm:"size:50" json:"reference_type"`
	ReferenceID     *uuid.UUID `gorm:"type:uuid" json:"reference_id"`
	Description     *string    `json:"description"`
	Notes           *string    `json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`

	Promoter User `gorm:"foreignKey:PromoterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

func (PromoterTransactionHistory) TableName() string {
	return "promoter_transaction_history"
}
