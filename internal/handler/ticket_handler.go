package handler

import (
	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	ticketUsecase usecase.TicketUsecase
}

func NewTicketHandler(ticketUsecase usecase.TicketUsecase) *TicketHandler {
	return &TicketHandler{ticketUsecase: ticketUsecase}
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
