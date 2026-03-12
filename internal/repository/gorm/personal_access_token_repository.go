package gorm

import (
	"time"

	"event-budaya-ticketing-bcc/internal/domain"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type personalAccessTokenRepository struct {
	db *gorm.DB
}

func NewPersonalAccessTokenRepository(db *gorm.DB) repository.PersonalAccessTokenRepository {
	return &personalAccessTokenRepository{db: db}
}

func (r *personalAccessTokenRepository) Create(token *domain.PersonalAccessToken) error {
	return r.db.Create(token).Error
}

func (r *personalAccessTokenRepository) FindByToken(token string) (*domain.PersonalAccessToken, error) {
	var pat domain.PersonalAccessToken
	err := r.db.Preload("User").Where("token = ?", token).First(&pat).Error
	if err != nil {
		return nil, err
	}
	return &pat, nil
}

func (r *personalAccessTokenRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&domain.PersonalAccessToken{}).Error
}

func (r *personalAccessTokenRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.PersonalAccessToken{}).Error
}

func (r *personalAccessTokenRepository) UpdateLastUsed(id uint, t time.Time) error {
	return r.db.Model(&domain.PersonalAccessToken{}).Where("id = ?", id).Update("last_used_at", t).Error
}
