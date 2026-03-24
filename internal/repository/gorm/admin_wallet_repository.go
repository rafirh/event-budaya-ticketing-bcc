package gorm

import (
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type adminWalletRepository struct {
	db *gorm.DB
}

func NewAdminWalletRepository(db *gorm.DB) repository.AdminWalletRepository {
	return &adminWalletRepository{db: db}
}

func (r *adminWalletRepository) FindOrCreate() (*model.AdminWallet, error) {
	var wallet model.AdminWallet
	err := r.db.First(&wallet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			wallet = model.AdminWallet{
				ID:      uuid.New(),
				Balance: 0,
			}
			if err := r.db.Create(&wallet).Error; err != nil {
				return nil, err
			}
			return &wallet, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *adminWalletRepository) Update(wallet *model.AdminWallet) error {
	return r.db.Save(wallet).Error
}

func (r *adminWalletRepository) FindByID(id string) (*model.AdminWallet, error) {
	var wallet model.AdminWallet
	err := r.db.Where("id = ?", id).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}
