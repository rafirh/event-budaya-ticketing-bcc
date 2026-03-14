package router

import (
	"event-budaya-ticketing-bcc/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func eventRoutes(api fiber.Router, eventHandler *handler.EventHandler) {
	events := api.Group("/events")
	events.Get("/", eventHandler.GetAll)
	events.Get("/:slug", eventHandler.GetBySlug)
}
