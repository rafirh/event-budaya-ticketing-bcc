package router

import (
	"event-budaya-ticketing-bcc/internal/handler"
	"event-budaya-ticketing-bcc/internal/middleware"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func orderRoutes(api fiber.Router, orderHandler *handler.OrderHandler, tokenRepo repository.PersonalAccessTokenRepository) {
	checkout := api.Group("/checkout", middleware.AuthMiddleware(tokenRepo))
	checkout.Post("", orderHandler.CreateTicketOrder)
}

func webhookRoutes(api fiber.Router, orderHandler *handler.OrderHandler) {
	webhook := api.Group("/webhook")
	webhook.Post("/midtrans", orderHandler.MidtransWebhook)
}
