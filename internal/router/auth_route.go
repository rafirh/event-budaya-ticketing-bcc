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
	auth.Post("/logout", middleware.AuthMiddleware(tokenRepo), authHandler.Logout)

	me := api.Group("/me", middleware.AuthMiddleware(tokenRepo))
	me.Get("", authHandler.Me)
	me.Patch("", authHandler.UpdateProfile)
}
