package router

import (
	"event-budaya-ticketing-bcc/internal/handler"
	"event-budaya-ticketing-bcc/internal/middleware"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func eventRoutes(api fiber.Router, eventHandler *handler.EventHandler, tokenRepo repository.PersonalAccessTokenRepository) {
	events := api.Group("/events")
	events.Get("/", eventHandler.GetAll)
	events.Get("/:slug", eventHandler.GetBySlug)

	promoterEvents := events.Group("", middleware.AuthMiddleware(tokenRepo), middleware.PromoterMiddleware())
	promoterEvents.Post("/", eventHandler.CreateEvent)

	promoterMeEvents := api.Group("/me/events", middleware.AuthMiddleware(tokenRepo), middleware.PromoterMiddleware())
	promoterMeEvents.Get("", eventHandler.GetByPromoterID)
}
