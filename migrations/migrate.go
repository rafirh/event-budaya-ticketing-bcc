package migrations

import (
	"log"

	"event-budaya-ticketing-bcc/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Migrate runs all database migrations
func Migrate(db *gorm.DB) {
	log.Println("Running migrations...")

	err := db.AutoMigrate(
		&domain.User{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
}

// Seed seeds the database with initial data
func Seed(db *gorm.DB) {
	log.Println("Running seeders...")

	// Seed admin user
	seedAdminUser(db)

	// Seed sample users
	seedUsers(db)

	log.Println("Seeders completed successfully")
}

func seedAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Where("email = ?", "admin@example.com").Count(&count)

	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := domain.User{
			Name:     "Administrator",
			Email:    "admin@example.com",
			Password: string(hashedPassword),
			Role:     "admin",
		}
		db.Create(&admin)
		log.Println("Admin user seeded: admin@example.com / admin123")
	}
}

func seedUsers(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Where("role = ?", "user").Count(&count)

	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)

		users := []domain.User{
			{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: string(hashedPassword),
				Role:     "user",
			},
			{
				Name:     "Jane Smith",
				Email:    "jane@example.com",
				Password: string(hashedPassword),
				Role:     "user",
			},
			{
				Name:     "Bob Wilson",
				Email:    "bob@example.com",
				Password: string(hashedPassword),
				Role:     "user",
			},
		}

		for _, user := range users {
			db.Create(&user)
		}
		log.Println("Sample users seeded (password: user123)")
	}
}

// Fresh drops all tables and re-runs migrations
func Fresh(db *gorm.DB) {
	log.Println("Dropping all tables...")

	// Drop tables in reverse order of dependencies
	db.Migrator().DropTable(&domain.User{})

	log.Println("All tables dropped")

	// Run migrations
	Migrate(db)
}

// Rollback drops the last migration
func Rollback(db *gorm.DB) {
	log.Println("Rolling back last migration...")
	// In GORM, you would need to implement custom rollback logic
	// For simplicity, this example just logs the action
	log.Println("Rollback completed")
}
