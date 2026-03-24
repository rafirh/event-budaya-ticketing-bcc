package handler

import (
	"math"
	"strconv"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"
	"event-budaya-ticketing-bcc/pkg/validator"

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

	sortBy := c.Query("sort_by", "created_at")
	allowedSortFields := map[string]bool{
		"title":                 true,
		"start_date":            true,
		"end_date":              true,
		"registration_deadline": true,
		"price":                 true,
		"quota":                 true,
		"sold":                  true,
		"created_at":            true,
	}
	if !allowedSortFields[sortBy] {
		return response.Error(c, fiber.StatusBadRequest, "invalid sort_by parameter")
	}

	sortOrder := c.Query("sort_order", "desc")
	if sortOrder != "asc" && sortOrder != "desc" {
		return response.Error(c, fiber.StatusBadRequest, "invalid sort_order parameter")
	}

	events, total, err := h.eventUsecase.GetAll(search, categoryID, sortBy, sortOrder, page, limit)
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

func (h *EventHandler) GetByPromoterID(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return response.Error(c, fiber.StatusUnauthorized, "User ID not found in context")
	}

	promoterID, err := uuid.Parse(userID.(string))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
	}

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

	sortBy := c.Query("sort_by", "created_at")
	allowedSortFields := map[string]bool{
		"title":                 true,
		"start_date":            true,
		"end_date":              true,
		"registration_deadline": true,
		"price":                 true,
		"quota":                 true,
		"sold":                  true,
		"created_at":            true,
	}
	if !allowedSortFields[sortBy] {
		return response.Error(c, fiber.StatusBadRequest, "invalid sort_by parameter")
	}

	sortOrder := c.Query("sort_order", "desc")
	if sortOrder != "asc" && sortOrder != "desc" {
		return response.Error(c, fiber.StatusBadRequest, "invalid sort_order parameter")
	}

	events, total, err := h.eventUsecase.GetByPromoterID(promoterID, search, categoryID, sortBy, sortOrder, page, limit)
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

func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return response.Error(c, fiber.StatusUnauthorized, "User ID not found in context")
	}

	promoterID, err := uuid.Parse(userID.(string))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var req dto.CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return response.Error(c, fiber.StatusBadRequest, "Validation error: "+errors[0].Message)
	}

	result, err := h.eventUsecase.CreateEvent(req, promoterID)
	if err != nil {
		switch err.Error() {
		case "category not found":
			return response.Error(c, fiber.StatusBadRequest, "Category not found")
		case "title is required":
			return response.Error(c, fiber.StatusBadRequest, "Title is required")
		case "quota must be greater than 0":
			return response.Error(c, fiber.StatusBadRequest, "Quota must be greater than 0")
		default:
			return response.Error(c, fiber.StatusInternalServerError, err.Error())
		}
	}

	return response.Success(c, fiber.StatusCreated, "Event created successfully, please complete payment", result)
}
