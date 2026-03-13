package migrations

import (
	"log"

	"event-budaya-ticketing-bcc/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	log.Println("Running migrations...")

	if db.Dialector.Name() == "postgres" {
		if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
			log.Fatalf("Failed to enable uuid extension: %v", err)
		}
	}

	err := db.AutoMigrate(
		&model.User{},
		&model.PersonalAccessToken{},
		&model.EventCategory{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
}

func Seed(db *gorm.DB) {
	log.Println("Running seeders...")
	seedUsers(db)
	seedEventCategories(db)
	log.Println("Seeders completed successfully")
}

func seedUsers(db *gorm.DB) {
	users := []struct {
		Name     string
		Email    string
		Password string
		Role     string
	}{
		{"Administrator", "admin@gmail.com", "admin123", "admin"},
		{"Promotor", "promotor@gmail.com", "promotor123", "promotor"},
		{"User", "user@gmail.com", "user123", "user"},
	}

	for _, u := range users {
		var count int64
		db.Model(&model.User{}).Where("email = ?", u.Email).Count(&count)
		if count > 0 {
			continue
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		user := model.User{
			Name:     u.Name,
			Email:    u.Email,
			Password: string(hashedPassword),
			Role:     u.Role,
		}
		db.Create(&user)
		log.Printf("Seeded: %s / %s", u.Email, u.Password)
	}
}

func Fresh(db *gorm.DB) {
	log.Println("Dropping all tables...")
	db.Migrator().DropTable(&model.PersonalAccessToken{}, &model.EventCategory{}, &model.User{})
	log.Println("All tables dropped")
	Migrate(db)
}

func seedEventCategories(db *gorm.DB) {
	categories := []string{
		"Seminar",
		"Workshop",
		"Konferensi",
		"Talkshow",
		"Festival",
		"Kompetisi",
		"Exhibition",
		"Concert",
	}

	for _, name := range categories {
		var count int64
		db.Model(&model.EventCategory{}).Where("name = ?", name).Count(&count)
		if count > 0 {
			continue
		}

		category := model.EventCategory{Name: name}
		db.Create(&category)
		log.Printf("Seeded category: %s", name)
	}
}
