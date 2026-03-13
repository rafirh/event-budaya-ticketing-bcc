package gorm

import (
	"time"

	"event-budaya-ticketing-bcc/internal/model"
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

func (r *personalAccessTokenRepository) Create(token *model.PersonalAccessToken) error {
	return r.db.Create(token).Error
}

func (r *personalAccessTokenRepository) FindByToken(token string) (*model.PersonalAccessToken, error) {
	var pat model.PersonalAccessToken
	err := r.db.Preload("User").Where("token = ?", token).First(&pat).Error
	if err != nil {
		return nil, err
	}
	return &pat, nil
}

func (r *personalAccessTokenRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&model.PersonalAccessToken{}).Error
}

func (r *personalAccessTokenRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.PersonalAccessToken{}).Error
}

func (r *personalAccessTokenRepository) UpdateLastUsed(id uint, t time.Time) error {
	return r.db.Model(&model.PersonalAccessToken{}).Where("id = ?", id).Update("last_used_at", t).Error
}
