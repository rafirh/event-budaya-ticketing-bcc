package migrations

import (
	"log"

	"event-budaya-ticketing-bcc/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedUsers(db *gorm.DB) {
	users := []struct {
		Name     string
		Email    string
		Password string
		Role     string
	}{
		{"Administrator", "admin@gmail.com", "admin123", "admin"},
		{"Promotor Example", "promotor@gmail.com", "promotor123", "promotor"},
		{"User Example", "user@gmail.com", "user123", "user"},
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
