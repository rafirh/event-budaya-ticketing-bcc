package gorm

import (
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"gorm.io/gorm"
)

type feeRepository struct {
	db *gorm.DB
}

func NewFeeRepository(db *gorm.DB) repository.FeeRepository {
	return &feeRepository{db: db}
}

func (r *feeRepository) FindByType(feeType string) (*model.Fee, error) {
	var fee model.Fee
	err := r.db.Where("fee_type = ?", feeType).First(&fee).Error
	if err != nil {
		return nil, err
	}
	return &fee, nil
}

func (r *feeRepository) FindAll() ([]model.Fee, error) {
	var fees []model.Fee
	err := r.db.Find(&fees).Error
	if err != nil {
		return nil, err
	}
	return fees, nil
}
