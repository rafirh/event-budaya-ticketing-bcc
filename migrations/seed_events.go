package migrations

import (
	"log"
	"time"

	"event-budaya-ticketing-bcc/internal/model"

	"gorm.io/gorm"
)

func seedEvents(db *gorm.DB) {
	var promoter model.User
	if err := db.Where("role = ?", "promotor").Order("created_at ASC").First(&promoter).Error; err != nil {
		log.Printf("Skipping event seed: promoter user not found: %v", err)
		return
	}

	var categories []model.EventCategory
	if err := db.Order("name ASC").Find(&categories).Error; err != nil {
		log.Printf("Skipping event seed: failed to load categories: %v", err)
		return
	}
	if len(categories) == 0 {
		log.Println("Skipping event seed: no categories found")
		return
	}

	eventSeeds := []struct {
		Title       string
		Description string
		Venue       string
		Address     string
		MapsURL     string
		StartDate   time.Time
		EndDate     time.Time
		IsPaid      bool
		Status      string
	}{
		{"Seminar Budaya Nusantara 2026", "Sesi edukatif tentang kekayaan budaya Nusantara.", "Gedung Kesenian Malang", "Jl. Soekarno Hatta No. 1, Malang", "https://maps.google.com/?q=Gedung+Kesenian+Malang", time.Date(2026, time.March, 21, 9, 0, 0, 0, time.UTC), time.Date(2026, time.March, 21, 12, 0, 0, 0, time.UTC), true, "published"},
		{"Workshop Batik Modern", "Praktik langsung membuat motif batik modern.", "Balai Kreatif Malang", "Jl. Ijen No. 12, Malang", "https://maps.google.com/?q=Balai+Kreatif+Malang", time.Date(2026, time.March, 22, 13, 0, 0, 0, time.UTC), time.Date(2026, time.March, 22, 16, 0, 0, 0, time.UTC), true, "published"},
		{"Konferensi Ekonomi Kreatif", "Konferensi untuk pelaku ekonomi kreatif dan budaya.", "Convention Hall BCC", "Jl. Veteran No. 8, Malang", "https://maps.google.com/?q=Convention+Hall+BCC", time.Date(2026, time.March, 23, 8, 30, 0, 0, time.UTC), time.Date(2026, time.March, 23, 17, 0, 0, 0, time.UTC), true, "draft"},
		{"Talkshow Komunitas Seni Kota", "Diskusi santai bersama pelaku seni lokal.", "Amphitheater Kota", "Jl. Tugu No. 5, Malang", "https://maps.google.com/?q=Amphitheater+Kota+Malang", time.Date(2026, time.March, 24, 18, 30, 0, 0, time.UTC), time.Date(2026, time.March, 24, 20, 30, 0, 0, time.UTC), false, "published"},
		{"Festival Kuliner Tradisional", "Perayaan budaya lewat ragam kuliner tradisional.", "Alun-Alun Malang", "Jl. Merdeka Selatan, Malang", "https://maps.google.com/?q=Alun-Alun+Malang", time.Date(2026, time.March, 25, 10, 0, 0, 0, time.UTC), time.Date(2026, time.March, 25, 21, 0, 0, 0, time.UTC), false, "published"},
		{"Kompetisi Tari Pelajar", "Ajang perlombaan tari tradisional antar pelajar.", "Gor Ken Arok", "Jl. Mayjen Sungkono No. 10, Malang", "https://maps.google.com/?q=Gor+Ken+Arok+Malang", time.Date(2026, time.March, 26, 9, 0, 0, 0, time.UTC), time.Date(2026, time.March, 26, 15, 0, 0, 0, time.UTC), true, "draft"},
		{"Exhibition Karya UMKM", "Pameran karya kreatif dan produk UMKM lokal.", "Malang Creative Center", "Jl. A. Yani No. 53, Malang", "https://maps.google.com/?q=Malang+Creative+Center", time.Date(2026, time.March, 27, 10, 0, 0, 0, time.UTC), time.Date(2026, time.March, 27, 19, 0, 0, 0, time.UTC), false, "published"},
		{"Concert Harmoni Budaya", "Pertunjukan musik dengan sentuhan budaya tradisional.", "Open Air Theater", "Jl. Jakarta No. 40, Malang", "https://maps.google.com/?q=Open+Air+Theater+Malang", time.Date(2026, time.March, 28, 19, 0, 0, 0, time.UTC), time.Date(2026, time.March, 28, 22, 0, 0, 0, time.UTC), true, "published"},
		{"Seminar Pariwisata Berkelanjutan", "Pembahasan strategi pariwisata budaya yang berkelanjutan.", "Hotel Tugu Malang", "Jl. Tugu No. 3, Malang", "https://maps.google.com/?q=Hotel+Tugu+Malang", time.Date(2026, time.March, 29, 9, 0, 0, 0, time.UTC), time.Date(2026, time.March, 29, 12, 30, 0, 0, time.UTC), true, "draft"},
		{"Workshop Fotografi Event", "Pelatihan teknik fotografi untuk dokumentasi event budaya.", "Studio Kreatif Arema", "Jl. Bandung No. 9, Malang", "https://maps.google.com/?q=Studio+Kreatif+Arema", time.Date(2026, time.March, 30, 13, 0, 0, 0, time.UTC), time.Date(2026, time.March, 30, 17, 0, 0, 0, time.UTC), true, "published"},
	}

	for index, seed := range eventSeeds {
		var count int64
		db.Model(&model.Event{}).Where("title = ?", seed.Title).Count(&count)
		if count > 0 {
			continue
		}

		slug := makeSlug(seed.Title)
		category := categories[index%len(categories)]
		description := seed.Description
		venue := seed.Venue
		address := seed.Address
		mapsURL := seed.MapsURL

		event := model.Event{
			PromoterID:    promoter.ID,
			CategoryID:    &category.ID,
			Title:         seed.Title,
			Slug:          &slug,
			Description:   &description,
			Venue:         &venue,
			Address:       &address,
			GoogleMapsURL: &mapsURL,
			StartDate:     &seed.StartDate,
			EndDate:       &seed.EndDate,
			IsPaid:        seed.IsPaid,
			BannerURL:     nil,
			Status:        seed.Status,
		}

		if err := db.Create(&event).Error; err != nil {
			log.Printf("Failed seeding event %s: %v", seed.Title, err)
			continue
		}

		log.Printf("Seeded event: %s", seed.Title)
	}
}
