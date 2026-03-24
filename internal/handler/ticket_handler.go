package handler

import (
	"math"
	"strconv"

	"event-budaya-ticketing-bcc/internal/repository"
	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TicketHandler struct {
	ticketUsecase   usecase.TicketUsecase
	eventRepository repository.EventRepository
}

func NewTicketHandler(ticketUsecase usecase.TicketUsecase, eventRepository repository.EventRepository) *TicketHandler {
	return &TicketHandler{
		ticketUsecase:   ticketUsecase,
		eventRepository: eventRepository,
	}
}

func (h *TicketHandler) GetMyTickets(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	tickets, err := h.ticketUsecase.GetMyTickets(userID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Tickets retrieved successfully", tickets)
}

func (h *TicketHandler) GetMyTicketDetail(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	ticketID := c.Params("id")
	if ticketID == "" {
		return response.Error(c, fiber.StatusBadRequest, "Ticket ID is required")
	}

	ticketDetail, err := h.ticketUsecase.GetMyTicketDetail(userID, ticketID)
	if err != nil {
		switch err.Error() {
		case "ticket not found":
			return response.Error(c, fiber.StatusNotFound, err.Error())
		case "unauthorized":
			return response.Error(c, fiber.StatusForbidden, err.Error())
		default:
			return response.Error(c, fiber.StatusBadRequest, err.Error())
		}
	}

	return response.Success(c, fiber.StatusOK, "Ticket detail retrieved successfully", ticketDetail)
}

func (h *TicketHandler) GetAttendeesByEventID(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return response.Error(c, fiber.StatusUnauthorized, "User ID not found in context")
	}

	promoterID, err := uuid.Parse(userID.(string))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	eventID := c.Params("eventId")
	if eventID == "" {
		return response.Error(c, fiber.StatusBadRequest, "Event ID is required")
	}

	// Verify event ownership
	event, err := h.eventRepository.FindByID(eventID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Event not found")
	}

	if event.PromoterID != promoterID {
		return response.Error(c, fiber.StatusForbidden, "You don't have permission to access this event's attendees")
	}

	search := c.Query("search")

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

	attendees, total, err := h.ticketUsecase.GetAttendeesByEventID(eventID, search, page, limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return response.Paginated(c, fiber.StatusOK, "Attendees retrieved successfully", attendees, response.Pagination{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	})
}
