package gorm

import (
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"gorm.io/gorm"
)

type promoterWalletRepository struct {
	db *gorm.DB
}

func NewPromoterWalletRepository(db *gorm.DB) repository.PromoterWalletRepository {
	return &promoterWalletRepository{db: db}
}

func (r *promoterWalletRepository) Create(wallet *model.PromotorWallet) error {
	return r.db.Create(wallet).Error
}

func (r *promoterWalletRepository) FindByPromoterID(promoterID string) (*model.PromotorWallet, error) {
	var wallet model.PromotorWallet
	err := r.db.Where("promoter_id = ?", promoterID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *promoterWalletRepository) Update(wallet *model.PromotorWallet) error {
	return r.db.Save(wallet).Error
}

type walletTransactionRepository struct {
	db *gorm.DB
}

func NewWalletTransactionRepository(db *gorm.DB) repository.WalletTransactionRepository {
	return &walletTransactionRepository{db: db}
}

func (r *walletTransactionRepository) Create(transaction *model.WalletTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *walletTransactionRepository) FindByWalletID(walletID string) ([]model.WalletTransaction, error) {
	var transactions []model.WalletTransaction
	err := r.db.Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
