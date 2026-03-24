package gorm

import (
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type promoterTransactionHistoryRepository struct {
	db *gorm.DB
}

func NewPromoterTransactionHistoryRepository(db *gorm.DB) repository.PromoterTransactionHistoryRepository {
	return &promoterTransactionHistoryRepository{db: db}
}

func (r *promoterTransactionHistoryRepository) Create(transaction *model.PromoterTransactionHistory) error {
	if transaction.ID == uuid.Nil {
		transaction.ID = uuid.New()
	}
	return r.db.Create(transaction).Error
}

func (r *promoterTransactionHistoryRepository) FindByPromoterID(promoterID uuid.UUID, limit, offset int) ([]model.PromoterTransactionHistory, int64, error) {
	var transactions []model.PromoterTransactionHistory
	var total int64

	if err := r.db.Model(&model.PromoterTransactionHistory{}).
		Where("promoter_id = ?", promoterID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Where("promoter_id = ?", promoterID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error

	return transactions, total, err
}
