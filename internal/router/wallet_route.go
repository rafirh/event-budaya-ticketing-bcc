package router

import (
	"event-budaya-ticketing-bcc/internal/handler"
	"event-budaya-ticketing-bcc/internal/middleware"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func SetupWalletRoute(api fiber.Router, walletHandler *handler.WalletHandler, tokenRepo repository.PersonalAccessTokenRepository) {
	wallet := api.Group("/me/wallet", middleware.AuthMiddleware(tokenRepo))
	wallet.Use(middleware.PromoterMiddleware())

	wallet.Get("", walletHandler.GetMyWallet)
}
