package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminWallet struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Balance        float64    `gorm:"type:numeric(14,2);default:0" json:"balance"`
	TotalRevenue   float64    `gorm:"type:numeric(14,2);default:0" json:"total_revenue"`
	TotalWithdrawn float64    `gorm:"type:numeric(14,2);default:0" json:"total_withdrawn"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

func (AdminWallet) TableName() string {
	return "admin_wallets"
}

func (a *AdminWallet) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
