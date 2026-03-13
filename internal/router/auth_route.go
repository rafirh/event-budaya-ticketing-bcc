package router

import (
	"event-budaya-ticketing-bcc/internal/handler"
	"event-budaya-ticketing-bcc/internal/middleware"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func authRoutes(api fiber.Router, authHandler *handler.AuthHandler, tokenRepo repository.PersonalAccessTokenRepository) {
	auth := api.Group("/auth")

	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	protected := auth.Group("", middleware.AuthMiddleware(tokenRepo))
	protected.Get("/me", authHandler.Me)
	protected.Patch("/profile", authHandler.UpdateProfile)
	protected.Post("/logout", authHandler.Logout)
}
