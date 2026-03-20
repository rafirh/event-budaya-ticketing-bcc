package router

import (
	"event-budaya-ticketing-bcc/internal/handler"
	"event-budaya-ticketing-bcc/internal/middleware"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func commentRoutes(api fiber.Router, eventCommentHandler *handler.EventCommentHandler, tokenRepo repository.PersonalAccessTokenRepository) {
	events := api.Group("/events")
	events.Get("/:eventId/comments", eventCommentHandler.GetByEventID)
	events.Post("/:eventId/comments", middleware.AuthMiddleware(tokenRepo), eventCommentHandler.Create)
}
