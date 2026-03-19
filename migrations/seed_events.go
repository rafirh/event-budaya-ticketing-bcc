package migrations

import (
	"log"
	"time"

	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/pkg/helper"

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
		Title                string
		Summary              string
		Description          string
		Venue                string
		Address              string
		MapsURL              string
		StartDate            time.Time
		EndDate              time.Time
		RegistrationDeadline time.Time
		IsPaid               bool
		Price                float64
		Quota                int
		Sold                 int
		Status               string
	}{
		{"Seminar Budaya Nusantara 2026", "Seminar inspiratif tentang kekayaan budaya Nusantara.", "Sesi edukatif tentang kekayaan budaya Nusantara.", "Gedung Kesenian Malang", "Jl. Soekarno Hatta No. 1, Malang", "https://maps.google.com/?q=Gedung+Kesenian+Malang", time.Date(2026, time.March, 21, 9, 0, 0, 0, time.UTC), time.Date(2026, time.March, 21, 12, 0, 0, 0, time.UTC), time.Date(2026, time.March, 20, 23, 59, 0, 0, time.UTC), true, 75000, 300, 120, "published"},
		{"Workshop Batik Modern", "Workshop praktis membuat motif batik modern.", "Praktik langsung membuat motif batik modern.", "Balai Kreatif Malang", "Jl. Ijen No. 12, Malang", "https://maps.google.com/?q=Balai+Kreatif+Malang", time.Date(2026, time.March, 22, 13, 0, 0, 0, time.UTC), time.Date(2026, time.March, 22, 16, 0, 0, 0, time.UTC), time.Date(2026, time.March, 21, 23, 59, 0, 0, time.UTC), true, 100000, 120, 35, "published"},
		{"Konferensi Ekonomi Kreatif", "Forum kolaborasi pelaku ekonomi kreatif budaya.", "Konferensi untuk pelaku ekonomi kreatif dan budaya.", "Convention Hall BCC", "Jl. Veteran No. 8, Malang", "https://maps.google.com/?q=Convention+Hall+BCC", time.Date(2026, time.March, 23, 8, 30, 0, 0, time.UTC), time.Date(2026, time.March, 23, 17, 0, 0, 0, time.UTC), time.Date(2026, time.March, 22, 23, 59, 0, 0, time.UTC), true, 150000, 500, 0, "draft"},
		{"Talkshow Komunitas Seni Kota", "Obrolan santai bersama pegiat seni lokal.", "Diskusi santai bersama pelaku seni lokal.", "Amphitheater Kota", "Jl. Tugu No. 5, Malang", "https://maps.google.com/?q=Amphitheater+Kota+Malang", time.Date(2026, time.March, 24, 18, 30, 0, 0, time.UTC), time.Date(2026, time.March, 24, 20, 30, 0, 0, time.UTC), time.Date(2026, time.March, 24, 12, 0, 0, 0, time.UTC), false, 0, 1000, 240, "published"},
		{"Festival Kuliner Tradisional", "Jelajah rasa khas kuliner tradisional Nusantara.", "Perayaan budaya lewat ragam kuliner tradisional.", "Alun-Alun Malang", "Jl. Merdeka Selatan, Malang", "https://maps.google.com/?q=Alun-Alun+Malang", time.Date(2026, time.March, 25, 10, 0, 0, 0, time.UTC), time.Date(2026, time.March, 25, 21, 0, 0, 0, time.UTC), time.Date(2026, time.March, 24, 23, 59, 0, 0, time.UTC), false, 0, 2000, 500, "published"},
		{"Kompetisi Tari Pelajar", "Kompetisi tari tradisional antar sekolah.", "Ajang perlombaan tari tradisional antar pelajar.", "Gor Ken Arok", "Jl. Mayjen Sungkono No. 10, Malang", "https://maps.google.com/?q=Gor+Ken+Arok+Malang", time.Date(2026, time.March, 26, 9, 0, 0, 0, time.UTC), time.Date(2026, time.March, 26, 15, 0, 0, 0, time.UTC), time.Date(2026, time.March, 25, 23, 59, 0, 0, time.UTC), true, 50000, 400, 0, "draft"},
		{"Exhibition Karya UMKM", "Pameran karya dan produk unggulan UMKM lokal.", "Pameran karya kreatif dan produk UMKM lokal.", "Malang Creative Center", "Jl. A. Yani No. 53, Malang", "https://maps.google.com/?q=Malang+Creative+Center", time.Date(2026, time.March, 27, 10, 0, 0, 0, time.UTC), time.Date(2026, time.March, 27, 19, 0, 0, 0, time.UTC), time.Date(2026, time.March, 26, 23, 59, 0, 0, time.UTC), false, 0, 800, 210, "published"},
		{"Concert Harmoni Budaya", "Konser musik modern dengan nuansa budaya tradisional.", "Pertunjukan musik dengan sentuhan budaya tradisional.", "Open Air Theater", "Jl. Jakarta No. 40, Malang", "https://maps.google.com/?q=Open+Air+Theater+Malang", time.Date(2026, time.March, 28, 19, 0, 0, 0, time.UTC), time.Date(2026, time.March, 28, 22, 0, 0, 0, time.UTC), time.Date(2026, time.March, 27, 23, 59, 0, 0, time.UTC), true, 120000, 600, 460, "published"},
		{"Seminar Pariwisata Berkelanjutan", "Strategi membangun pariwisata budaya berkelanjutan.", "Pembahasan strategi pariwisata budaya yang berkelanjutan.", "Hotel Tugu Malang", "Jl. Tugu No. 3, Malang", "https://maps.google.com/?q=Hotel+Tugu+Malang", time.Date(2026, time.March, 29, 9, 0, 0, 0, time.UTC), time.Date(2026, time.March, 29, 12, 30, 0, 0, time.UTC), time.Date(2026, time.March, 28, 23, 59, 0, 0, time.UTC), true, 90000, 250, 0, "draft"},
		{"Workshop Fotografi Event", "Belajar teknik foto untuk dokumentasi event budaya.", "Pelatihan teknik fotografi untuk dokumentasi event budaya.", "Studio Kreatif Arema", "Jl. Bandung No. 9, Malang", "https://maps.google.com/?q=Studio+Kreatif+Arema", time.Date(2026, time.March, 30, 13, 0, 0, 0, time.UTC), time.Date(2026, time.March, 30, 17, 0, 0, 0, time.UTC), time.Date(2026, time.March, 29, 23, 59, 0, 0, time.UTC), true, 65000, 180, 60, "published"},
	}

	for index, seed := range eventSeeds {
		var count int64
		db.Model(&model.Event{}).Where("title = ?", seed.Title).Count(&count)
		if count > 0 {
			continue
		}

		slug := helper.MakeSlug(seed.Title)
		category := categories[index%len(categories)]
		summary := seed.Summary
		description := seed.Description
		venue := seed.Venue
		address := seed.Address
		mapsURL := seed.MapsURL
		registrationDeadline := seed.RegistrationDeadline

		event := model.Event{
			PromoterID:           promoter.ID,
			CategoryID:           &category.ID,
			Title:                seed.Title,
			Slug:                 &slug,
			Summary:              &summary,
			Description:          &description,
			Venue:                &venue,
			Address:              &address,
			GoogleMapsURL:        &mapsURL,
			StartDate:            &seed.StartDate,
			EndDate:              &seed.EndDate,
			RegistrationDeadline: &registrationDeadline,
			IsPaid:               seed.IsPaid,
			Price:                seed.Price,
			Quota:                seed.Quota,
			Sold:                 seed.Sold,
			BannerURL:            nil,
			Status:               seed.Status,
		}

		if err := db.Create(&event).Error; err != nil {
			log.Printf("Failed seeding event %s: %v", seed.Title, err)
			continue
		}

		log.Printf("Seeded event: %s", seed.Title)
	}
}
