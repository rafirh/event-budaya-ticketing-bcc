package handler

import (
	"math"
	"strconv"

	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type EventHandler struct {
	eventUsecase usecase.EventUsecase
}

func NewEventHandler(eventUsecase usecase.EventUsecase) *EventHandler {
	return &EventHandler{eventUsecase: eventUsecase}
}

func (h *EventHandler) GetAll(c *fiber.Ctx) error {
	search := c.Query("search")
	categoryID := c.Query("category_id")
	if categoryID == "" {
		categoryID = c.Query("category")
	}
	if categoryID != "" {
		if _, err := uuid.Parse(categoryID); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid category_id parameter")
		}
	}

	page := 1
	if pageParam := c.Query("page"); pageParam != "" {
		parsedPage, err := strconv.Atoi(pageParam)
		if err != nil || parsedPage < 1 {
			return response.Error(c, fiber.StatusBadRequest, "invalid page parameter")
		}
		page = parsedPage
	}

	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err != nil || parsedLimit < 1 {
			return response.Error(c, fiber.StatusBadRequest, "invalid limit parameter")
		}
		limit = parsedLimit
	}

	events, total, err := h.eventUsecase.GetAll(search, categoryID, page, limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return response.Paginated(c, fiber.StatusOK, "Events retrieved successfully", events, response.Pagination{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	})
}

func (h *EventHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	event, err := h.eventUsecase.GetBySlug(slug)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Event retrieved successfully", event)
}
