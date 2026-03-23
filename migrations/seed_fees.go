package migrations

import (
	"log"

	"event-budaya-ticketing-bcc/internal/model"

	"gorm.io/gorm"
)

func seedFees(db *gorm.DB) {
	fees := []struct {
		FeeType         string
		CalculationType string
		Amount          float64
	}{
		{
			FeeType:         "SERVICE_FEE",
			CalculationType: "fixed",
			Amount:          2000,
		},
		{
			FeeType:         "EVENT_POSTING_FEE",
			CalculationType: "fixed",
			Amount:          20000,
		},
	}

	for _, f := range fees {
		var count int64
		db.Model(&model.Fee{}).Where("fee_type = ?", f.FeeType).Count(&count)
		if count > 0 {
			continue
		}

		fee := model.Fee{
			FeeType:         f.FeeType,
			CalculationType: f.CalculationType,
			Amount:          f.Amount,
		}
		db.Create(&fee)
		log.Printf("Seeded: %s (%.2f)", f.FeeType, f.Amount)
	}
}
