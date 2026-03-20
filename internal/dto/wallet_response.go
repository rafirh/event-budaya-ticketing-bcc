package dto

import (
	"time"

	"github.com/google/uuid"
)

type PromoterWalletResponse struct {
	ID         uuid.UUID `json:"id"`
	PromoterID uuid.UUID `json:"promoter_id"`
	Balance    float64   `json:"balance"`
	CreatedAt  time.Time `json:"created_at"`
}

type WalletTransactionResponse struct {
	ID          uuid.UUID  `json:"id"`
	WalletID    uuid.UUID  `json:"wallet_id"`
	Type        string     `json:"type"`
	Direction   string     `json:"direction"`
	Amount      float64    `json:"amount"`
	ReferenceID *uuid.UUID `json:"reference_id"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
}
