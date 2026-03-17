package handler

import (
	"math"
	"strconv"

	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryUsecase usecase.CategoryUsecase
}

func NewCategoryHandler(categoryUsecase usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{categoryUsecase: categoryUsecase}
}

func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
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

	sortBy := c.Query("sort_by", "name")
	allowedSortFields := map[string]bool{
		"name":        true,
		"event_count": true,
	}
	if !allowedSortFields[sortBy] {
		return response.Error(c, fiber.StatusBadRequest, "invalid sort_by parameter")
	}

	sortOrder := c.Query("sort_order", "asc")
	if sortOrder != "asc" && sortOrder != "desc" {
		return response.Error(c, fiber.StatusBadRequest, "invalid sort_order parameter")
	}

	categories, total, err := h.categoryUsecase.GetAll(page, limit, sortBy, sortOrder)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return response.Paginated(c, fiber.StatusOK, "Categories retrieved successfully", categories, response.Pagination{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	})
}
