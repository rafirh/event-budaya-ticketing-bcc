package gorm

import (
	"time"

	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type emailVerificationTokenRepository struct {
	db *gorm.DB
}

func NewEmailVerificationTokenRepository(db *gorm.DB) repository.EmailVerificationTokenRepository {
	return &emailVerificationTokenRepository{db: db}
}

func (r *emailVerificationTokenRepository) Create(token *model.EmailVerificationToken) error {
	return r.db.Create(token).Error
}

func (r *emailVerificationTokenRepository) FindByTokenHash(tokenHash string) (*model.EmailVerificationToken, error) {
	var token model.EmailVerificationToken
	err := r.db.Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *emailVerificationTokenRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.EmailVerificationToken{}).Error
}

func (r *emailVerificationTokenRepository) MarkUsed(id uint, usedAt time.Time) error {
	return r.db.Model(&model.EmailVerificationToken{}).Where("id = ?", id).Update("used_at", usedAt).Error
}
