package migrations

import (
	"log"

	"event-budaya-ticketing-bcc/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedAdminWallet(db *gorm.DB) {
	var count int64
	db.Model(&model.AdminWallet{}).Count(&count)
	if count > 0 {
		return
	}

	adminWallet := model.AdminWallet{
		ID:             uuid.New(),
		Balance:        0,
		TotalRevenue:   0,
		TotalWithdrawn: 0,
	}
	db.Create(&adminWallet)
	log.Printf("Seeded: Admin wallet (ID: %s)", adminWallet.ID)
}
