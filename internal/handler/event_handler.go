package handler

import (
	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	eventUsecase usecase.EventUsecase
}

func NewEventHandler(eventUsecase usecase.EventUsecase) *EventHandler {
	return &EventHandler{eventUsecase: eventUsecase}
}

func (h *EventHandler) GetAll(c *fiber.Ctx) error {
	events, err := h.eventUsecase.GetAll()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Events retrieved successfully", events)
}

func (h *EventHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	event, err := h.eventUsecase.GetBySlug(slug)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Event retrieved successfully", event)
}
