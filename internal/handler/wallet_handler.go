package handler

import (
	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	walletUsecase usecase.WalletUsecase
}

func NewWalletHandler(walletUsecase usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{walletUsecase: walletUsecase}
}

func (h *WalletHandler) GetMyWallet(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	wallet, err := h.walletUsecase.GetMyWallet(userID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Wallet retrieved successfully", wallet)
}
