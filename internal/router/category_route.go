package router

import (
	"event-budaya-ticketing-bcc/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func categoryRoutes(api fiber.Router, categoryHandler *handler.CategoryHandler) {
	categories := api.Group("/categories")
	categories.Get("/", categoryHandler.GetAll)
}
