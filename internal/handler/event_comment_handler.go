package handler

import (
	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"
	"event-budaya-ticketing-bcc/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type EventCommentHandler struct {
	eventCommentUsecase usecase.EventCommentUsecase
}

func NewEventCommentHandler(eventCommentUsecase usecase.EventCommentUsecase) *EventCommentHandler {
	return &EventCommentHandler{eventCommentUsecase: eventCommentUsecase}
}

func (h *EventCommentHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	eventID := c.Params("eventId")
	if eventID == "" {
		return response.Error(c, fiber.StatusBadRequest, "Event ID is required")
	}

	var req dto.CreateEventCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return response.ValidationError(c, errors)
	}

	result, err := h.eventCommentUsecase.CreateComment(userID, eventID, &req)
	if err != nil {
		switch err.Error() {
		case "event not found", "parent comment not found":
			return response.Error(c, fiber.StatusNotFound, err.Error())
		default:
			return response.Error(c, fiber.StatusBadRequest, err.Error())
		}
	}

	return response.Success(c, fiber.StatusCreated, "Comment created successfully", result)
}

func (h *EventCommentHandler) GetByEventID(c *fiber.Ctx) error {
	eventID := c.Params("eventId")
	if eventID == "" {
		return response.Error(c, fiber.StatusBadRequest, "Event ID is required")
	}

	comments, err := h.eventCommentUsecase.GetByEventID(eventID)
	if err != nil {
		switch err.Error() {
		case "event not found":
			return response.Error(c, fiber.StatusNotFound, err.Error())
		default:
			return response.Error(c, fiber.StatusBadRequest, err.Error())
		}
	}

	return response.Success(c, fiber.StatusOK, "Comments retrieved successfully", comments)
}
