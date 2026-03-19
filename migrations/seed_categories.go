package migrations

import (
	"log"

	"event-budaya-ticketing-bcc/internal/model"

	"gorm.io/gorm"
)

func seedEventCategories(db *gorm.DB) {
	baseIconURL := "https://event-budaya.iccn.or.id/seeders/categories/"

	categories := []struct {
		Name     string
		IconFile string
	}{
		{"Seni Pertunjukan", "seni_pertunjukan.svg"},
		{"Musik Tradisional", "musik_tradisional.svg"},
		{"Tari Daerah", "tari_daerah.svg"},
		{"Pameran Seni", "pameran_seni.svg"},
		{"Upacara Adat", "upacara_adat.svg"},
		{"Festival Adat", "festival_adat.svg"},
		{"Kuliner Budaya", "kuliner_budaya.svg"},
		{"Seminar Budaya", "seminar_budaya.svg"},
	}

	for _, cat := range categories {
		var count int64
		db.Model(&model.EventCategory{}).Where("name = ?", cat.Name).Count(&count)
		if count > 0 {
			continue
		}

		IconURL := baseIconURL + cat.IconFile
		category := model.EventCategory{
			Name: cat.Name,
			Icon: &IconURL,
		}
		db.Create(&category)
		log.Printf("Seeded category: %s", cat.Name)
	}
}
