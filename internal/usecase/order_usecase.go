package usecase

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"
	"event-budaya-ticketing-bcc/pkg/payment"

	"github.com/google/uuid"
)

const serviceFeePerTicket = 2000.0

type OrderUsecase interface {
	CreateTicketOrder(userID string, req *dto.CreateTicketOrderRequest) (*dto.CreateTicketOrderResponse, error)
	HandleMidtransWebhook(req *dto.MidtransWebhookRequest) error
}

type orderUsecase struct {
	userRepo       repository.UserRepository
	eventRepo      repository.EventRepository
	orderRepo      repository.OrderRepository
	ticketRepo     repository.TicketRepository
	paymentRepo    repository.PaymentRepository
	midtransClient *payment.Client
	midtransServer string
}

func NewOrderUsecase(
	userRepo repository.UserRepository,
	eventRepo repository.EventRepository,
	orderRepo repository.OrderRepository,
	ticketRepo repository.TicketRepository,
	paymentRepo repository.PaymentRepository,
	midtransClient *payment.Client,
	midtransServer string,
) OrderUsecase {
	return &orderUsecase{
		userRepo:       userRepo,
		eventRepo:      eventRepo,
		orderRepo:      orderRepo,
		ticketRepo:     ticketRepo,
		paymentRepo:    paymentRepo,
		midtransClient: midtransClient,
		midtransServer: midtransServer,
	}
}

func (u *orderUsecase) CreateTicketOrder(userID string, req *dto.CreateTicketOrderRequest) (*dto.CreateTicketOrderResponse, error) {
	if u.midtransClient == nil || u.midtransServer == "" {
		return nil, errors.New("midtrans is not configured")
	}

	event, err := u.eventRepo.FindByID(req.EventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	ticketCount := len(req.Tickets)
	if ticketCount < 1 {
		return nil, errors.New("at least one ticket is required")
	}

	if event.Sold >= event.Quota {
		return nil, errors.New("ticket is sold out")
	}

	if (event.Sold + ticketCount) > event.Quota {
		remaining := event.Quota - event.Sold
		return nil, fmt.Errorf("only %d tickets available", remaining)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user session")
	}

	unitPrice := event.Price
	serviceFeeTotal := float64(ticketCount) * serviceFeePerTicket
	totalPrice := (unitPrice * float64(ticketCount)) + serviceFeeTotal

	order := &model.Order{
		UserID:     parsedUserID,
		EventID:    event.ID,
		Quantity:   ticketCount,
		UnitPrice:  unitPrice,
		ServiceFee: serviceFeeTotal,
		TotalPrice: totalPrice,
		Status:     "pending",
	}
	if err := u.orderRepo.Create(order); err != nil {
		return nil, errors.New("failed to create order")
	}

	tickets := make([]model.Ticket, 0, ticketCount)
	for _, t := range req.Tickets {
		notes := ""
		if t.Notes != nil {
			notes = strings.TrimSpace(*t.Notes)
		}

		ticket := model.Ticket{
			OrderID:        order.ID,
			TicketCode:     generateTicketCode(order.ID),
			HolderName:     t.HolderName,
			IdentityType:   t.IdentityType,
			IdentityNumber: t.IdentityNumber,
			HolderPhone:    t.HolderPhone,
			HolderEmail:    t.HolderEmail,
			Notes:          notes,
			IsUsed:         false,
		}
		tickets = append(tickets, ticket)
	}

	if err := u.ticketRepo.CreateBatch(tickets); err != nil {
		return nil, errors.New("failed to create tickets")
	}

	paymentData := &model.Payment{
		OrderID: order.ID,
		Amount:  totalPrice,
		Status:  "waiting",
	}
	gateway := "midtrans"
	paymentData.PaymentGateway = &gateway
	if err := u.paymentRepo.Create(paymentData); err != nil {
		return nil, errors.New("failed to create payment")
	}

	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	snapResp, err := u.midtransClient.CreateSnapTransaction(payment.SnapTransactionRequest{
		TransactionDetails: payment.TransactionDetails{
			OrderID:  order.ID.String(),
			GrossAmt: int64(totalPrice),
		},
		CustomerDetails: &payment.CustomerDetails{
			FirstName: user.Name,
			Email:     user.Email,
			Phone:     derefString(user.Phone),
		},
	})
	if err != nil {
		return nil, errors.New("failed to create midtrans transaction")
	}

	paymentURL := snapResp.RedirectURL
	paymentData.PaymentURL = &paymentURL
	if err := u.paymentRepo.Update(paymentData); err != nil {
		return nil, errors.New("failed to update payment url")
	}

	event.Sold += ticketCount
	if event.Sold > event.Quota {
		return nil, errors.New("failed to reserve event tickets")
	}
	if err := u.eventRepo.Update(event); err != nil {
		return nil, errors.New("failed to reserve event tickets")
	}

	return &dto.CreateTicketOrderResponse{
		OrderID:         order.ID,
		EventID:         order.EventID,
		TicketCount:     ticketCount,
		UnitPrice:       unitPrice,
		ServiceFee:      serviceFeePerTicket,
		ServiceFeeTotal: serviceFeeTotal,
		TotalPrice:      totalPrice,
		PaymentStatus:   paymentData.Status,
		PaymentToken:    snapResp.Token,
		PaymentURL:      snapResp.RedirectURL,
	}, nil
}

func (u *orderUsecase) HandleMidtransWebhook(req *dto.MidtransWebhookRequest) error {
	if u.midtransServer == "" {
		return errors.New("midtrans server key is not configured")
	}

	if req.OrderID == "" || req.StatusCode == "" || req.GrossAmount == "" || req.SignatureKey == "" {
		return errors.New("invalid midtrans payload")
	}

	if !isValidMidtransSignature(req, u.midtransServer) {
		return errors.New("invalid midtrans signature")
	}

	order, err := u.orderRepo.FindByID(req.OrderID)
	if err != nil {
		return errors.New("order not found")
	}

	paymentData, err := u.paymentRepo.FindByOrderID(req.OrderID)
	if err != nil {
		return errors.New("payment not found")
	}

	paymentStatus, orderStatus := mapMidtransStatus(req.TransactionStatus, req.FraudStatus)

	paymentData.Status = paymentStatus
	method := req.PaymentType
	if method != "" {
		paymentData.PaymentMethod = &method
	}
	if paymentStatus == "success" {
		now := time.Now()
		paymentData.PaidAt = &now
	}
	if err := u.paymentRepo.Update(paymentData); err != nil {
		return errors.New("failed to update payment")
	}

	wasPaid := order.Status == "paid"
	wasPending := order.Status == "pending"
	order.Status = orderStatus
	if err := u.orderRepo.Update(order); err != nil {
		return errors.New("failed to update order")
	}

	if orderStatus == "paid" && !wasPaid {
		return nil
	}

	if orderStatus == "cancelled" && wasPending {
		event, err := u.eventRepo.FindByID(order.EventID.String())
		if err != nil {
			return errors.New("event not found")
		}
		event.Sold -= order.Quantity
		if event.Sold < 0 {
			event.Sold = 0
		}
		if err := u.eventRepo.Update(event); err != nil {
			return errors.New("failed to restore event tickets")
		}
	}

	return nil
}

func mapMidtransStatus(transactionStatus, fraudStatus string) (paymentStatus, orderStatus string) {
	switch transactionStatus {
	case "settlement":
		return "success", "paid"
	case "capture":
		if fraudStatus == "accept" {
			return "success", "paid"
		}
		return "waiting", "pending"
	case "pending":
		return "waiting", "pending"
	case "deny", "cancel", "expire", "failure":
		return "failed", "cancelled"
	default:
		return "waiting", "pending"
	}
}

func isValidMidtransSignature(req *dto.MidtransWebhookRequest, serverKey string) bool {
	raw := req.OrderID + req.StatusCode + req.GrossAmount + serverKey
	sum := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(sum[:])
	return strings.EqualFold(expected, req.SignatureKey)
}

func generateTicketCode(orderID uuid.UUID) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("TIX-%s-%d", strings.ToUpper(orderID.String()[0:8]), timestamp)
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
