package migrations

import (
	"log"

	"event-budaya-ticketing-bcc/internal/model"

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
		&model.Event{},
		&model.Order{},
		&model.Ticket{},
		&model.Payment{},
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
	seedEvents(db)
	log.Println("Seeders completed successfully")
}

func Fresh(db *gorm.DB) {
	log.Println("Dropping all tables...")
	db.Migrator().DropTable(&model.Payment{}, &model.Ticket{}, &model.Order{}, &model.PersonalAccessToken{}, &model.Event{}, &model.EventCategory{}, &model.User{})
	log.Println("All tables dropped")
	Migrate(db)
}
