package handler

import (
	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"
	"event-budaya-ticketing-bcc/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUsecase: orderUsecase}
}

func (h *OrderHandler) CreateTicketOrder(c *fiber.Ctx) error {
	var req dto.CreateTicketOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return response.ValidationError(c, errors)
	}
	for _, ticket := range req.Tickets {
		if errors := validator.ValidateStruct(ticket); errors != nil {
			return response.ValidationError(c, errors)
		}
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	result, err := h.orderUsecase.CreateTicketOrder(userID, &req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Ticket order created successfully", result)
}

func (h *OrderHandler) MidtransWebhook(c *fiber.Ctx) error {
	var req dto.MidtransWebhookRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid webhook payload")
	}

	if err := h.orderUsecase.HandleMidtransWebhook(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Webhook processed", nil)
}

func (h *OrderHandler) GetMyOrders(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	orders, err := h.orderUsecase.GetMyOrders(userID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Orders retrieved successfully", orders)
}

func (h *OrderHandler) GetMyOrderDetail(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	orderID := c.Params("id")
	if orderID == "" {
		return response.Error(c, fiber.StatusBadRequest, "Order ID is required")
	}

	orderDetail, err := h.orderUsecase.GetMyOrderDetail(userID, orderID)
	if err != nil {
		switch err.Error() {
		case "order not found":
			return response.Error(c, fiber.StatusNotFound, err.Error())
		case "unauthorized":
			return response.Error(c, fiber.StatusForbidden, err.Error())
		default:
			return response.Error(c, fiber.StatusBadRequest, err.Error())
		}
	}

	return response.Success(c, fiber.StatusOK, "Order detail retrieved successfully", orderDetail)
}
