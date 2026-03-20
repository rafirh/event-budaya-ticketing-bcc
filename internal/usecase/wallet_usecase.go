package usecase

import (
	"errors"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/google/uuid"
)

type WalletUsecase interface {
	GetMyWallet(promoterID string) (*dto.PromoterWalletResponse, error)
}

type walletUsecase struct {
	walletRepo repository.PromoterWalletRepository
}

func NewWalletUsecase(
	walletRepo repository.PromoterWalletRepository,
) WalletUsecase {
	return &walletUsecase{
		walletRepo: walletRepo,
	}
}

func (u *walletUsecase) GetMyWallet(promoterID string) (*dto.PromoterWalletResponse, error) {
	if promoterID == "" {
		return nil, errors.New("invalid promoter id")
	}

	parsedID, err := uuid.Parse(promoterID)
	if err != nil {
		return nil, errors.New("invalid promoter id format")
	}

	wallet, err := u.walletRepo.FindByPromoterID(promoterID)
	if err != nil {
		wallet = &model.PromotorWallet{
			PromoterID: parsedID,
			Balance:    0,
		}
		if err := u.walletRepo.Create(wallet); err != nil {
			return nil, errors.New("failed to create wallet")
		}
	}

	return &dto.PromoterWalletResponse{
		ID:         wallet.ID,
		PromoterID: wallet.PromoterID,
		Balance:    wallet.Balance,
		CreatedAt:  wallet.CreatedAt,
	}, nil
}
