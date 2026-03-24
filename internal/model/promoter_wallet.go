package model

import (
	"time"

	"github.com/google/uuid"
)

type PromotorWallet struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PromoterID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"promoter_id"`
	Balance    float64   `gorm:"type:numeric(14,2);default:0" json:"balance"`
	CreatedAt  time.Time `json:"created_at"`
	Promoter   User      `gorm:"foreignKey:PromoterID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

func (PromotorWallet) TableName() string {
	return "promoter_wallets"
}

type WalletTransaction struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	WalletID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"wallet_id"`
	Type        string         `gorm:"size:50;not null" json:"type"`
	Direction   string         `gorm:"size:10;not null" json:"direction"`
	Amount      float64        `gorm:"type:numeric(12,2);not null" json:"amount"`
	ReferenceID *uuid.UUID     `gorm:"type:uuid" json:"reference_id"`
	Description *string        `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	Wallet      PromotorWallet `gorm:"foreignKey:WalletID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}
