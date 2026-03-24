package usecase

import (
	"errors"
	"fmt"
	"time"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"
	"event-budaya-ticketing-bcc/pkg/helper"
	"event-budaya-ticketing-bcc/pkg/payment"

	"github.com/google/uuid"
)

type EventUsecase interface {
	GetAll(search, categoryID, sortBy, sortOrder string, page, limit int) ([]dto.EventListResponse, int64, error)
	GetBySlug(slug string) (*dto.EventDetailResponse, error)
	CreateEvent(req dto.CreateEventRequest, promoterID uuid.UUID) (*dto.EventCreatedWithPayment, error)
	HandleEventPaymentWebhook(req *dto.MidtransWebhookRequest) error
}

type eventUsecase struct {
	eventRepo                repository.EventRepository
	eventCreationPaymentRepo repository.EventCreationPaymentRepository
	feeRepo                  repository.FeeRepository
	adminWalletRepo          repository.AdminWalletRepository
	promoterTransactionRepo  repository.PromoterTransactionHistoryRepository
	paymentGateway           *payment.Client
}

func NewEventUsecase(
	eventRepo repository.EventRepository,
	eventCreationPaymentRepo repository.EventCreationPaymentRepository,
	feeRepo repository.FeeRepository,
	adminWalletRepo repository.AdminWalletRepository,
	promoterTransactionRepo repository.PromoterTransactionHistoryRepository,
	paymentGateway *payment.Client,
) EventUsecase {
	return &eventUsecase{
		eventRepo:                eventRepo,
		eventCreationPaymentRepo: eventCreationPaymentRepo,
		feeRepo:                  feeRepo,
		adminWalletRepo:          adminWalletRepo,
		promoterTransactionRepo:  promoterTransactionRepo,
		paymentGateway:           paymentGateway,
	}
}

func (u *eventUsecase) GetAll(search, categoryID, sortBy, sortOrder string, page, limit int) ([]dto.EventListResponse, int64, error) {
	offset := (page - 1) * limit

	events, total, err := u.eventRepo.FindAll(search, categoryID, sortBy, sortOrder, limit, offset)
	if err != nil {
		return nil, 0, errors.New("failed to fetch events")
	}

	responses := make([]dto.EventListResponse, 0, len(events))
	for _, event := range events {
		responses = append(responses, dto.ToEventListResponse(event))
	}

	return responses, total, nil
}

func (u *eventUsecase) GetBySlug(slug string) (*dto.EventDetailResponse, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}

	event, err := u.eventRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("event not found")
	}

	response := dto.ToEventDetailResponse(*event)
	return &response, nil
}

func (u *eventUsecase) CreateEvent(req dto.CreateEventRequest, promoterID uuid.UUID) (*dto.EventCreatedWithPayment, error) {
	// Validate request
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.Quota <= 0 {
		return nil, errors.New("quota must be greater than 0")
	}

	// Get fee for event creation
	fee, err := u.feeRepo.FindByType("EVENT_POSTING_FEE")
	if err != nil {
		return nil, errors.New("fee setting not found")
	}

	// Create slug
	slug := helper.MakeSlug(req.Title)

	// Create event model
	event := model.Event{
		ID:                   uuid.New(),
		PromoterID:           promoterID,
		CategoryID:           req.CategoryID,
		Title:                req.Title,
		Slug:                 &slug,
		Summary:              req.Summary,
		Description:          req.Description,
		Venue:                req.Venue,
		Address:              req.Address,
		GoogleMapsURL:        req.GoogleMapsURL,
		StartDate:            req.StartDate,
		EndDate:              req.EndDate,
		RegistrationDeadline: req.RegistrationDeadline,
		IsPaid:               req.IsPaid,
		Price:                req.Price,
		Quota:                req.Quota,
		BannerURL:            req.BannerURL,
		Status:               "draft", // Draft hingga pembayaran selesai
		CreatedAt:            time.Now(),
	}

	// Save event to database
	if err := u.eventRepo.Create(&event); err != nil {
		return nil, errors.New("failed to create event")
	}

	// Create payment via Midtrans
	orderID := fmt.Sprintf("EVT-%s-%d", event.ID.String()[:8], time.Now().Unix())

	snapReq := payment.SnapTransactionRequest{
		TransactionDetails: payment.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(fee.Amount),
		},
	}

	snapResp, err := u.paymentGateway.CreateSnapTransaction(snapReq)
	if err != nil {
		return nil, errors.New("failed to create payment transaction")
	}

	// Save payment record
	eventPayment := model.EventCreationPayment{
		ID:           uuid.New(),
		EventID:      event.ID,
		PromoterID:   promoterID,
		OrderID:      orderID,
		PaymentToken: snapResp.Token,
		Amount:       fee.Amount,
		Status:       "pending",
		PaymentURL:   &snapResp.RedirectURL,
		CreatedAt:    time.Now(),
	}

	if err := u.eventCreationPaymentRepo.Create(&eventPayment); err != nil {
		return nil, errors.New("failed to save payment record")
	}

	// Return response with payment link
	expiresAt := time.Now().Add(1 * time.Hour)

	response := &dto.EventCreatedWithPayment{
		EventID:      event.ID,
		Title:        event.Title,
		PaymentToken: snapResp.Token,
		PaymentURL:   snapResp.RedirectURL,
		Amount:       fee.Amount,
		Description:  fmt.Sprintf("Event creation fee for: %s", event.Title),
		ExpiresAt:    expiresAt,
	}

	return response, nil
}

func (u *eventUsecase) HandleEventPaymentWebhook(req *dto.MidtransWebhookRequest) error {
	// Find payment record by order ID
	payment, err := u.eventCreationPaymentRepo.FindByOrderID(req.OrderID)
	if err != nil {
		return errors.New("payment record not found for order ID: " + req.OrderID)
	}

	// Update payment status based on transaction status
	paymentStatus := "failure"
	if req.TransactionStatus == "settlement" || req.TransactionStatus == "capture" {
		paymentStatus = "settlement"
	} else if req.TransactionStatus == "pending" {
		paymentStatus = "pending"
		return nil // Don't process yet
	} else if req.TransactionStatus == "expire" || req.TransactionStatus == "cancel" {
		paymentStatus = "cancelled"
	}

	payment.Status = paymentStatus
	if err := u.eventCreationPaymentRepo.Update(payment); err != nil {
		return errors.New("failed to update payment status")
	}

	// If payment failed or cancelled, don't process further
	if paymentStatus != "settlement" {
		return nil
	}

	// Get event
	event, err := u.eventRepo.FindByID(payment.EventID.String())
	if err != nil {
		return errors.New("event not found")
	}

	// Update event status to published
	event.Status = "published"
	if err := u.eventRepo.Update(event); err != nil {
		return errors.New("failed to publish event")
	}

	// Get admin wallet (create if not exists)
	adminWallet, err := u.adminWalletRepo.FindOrCreate()
	if err != nil {
		return errors.New("failed to get admin wallet")
	}

	// Update admin wallet balance and revenue
	adminWallet.Balance += payment.Amount
	adminWallet.TotalRevenue += payment.Amount
	if err := u.adminWalletRepo.Update(adminWallet); err != nil {
		return errors.New("failed to update admin wallet")
	}

	// Create transaction record in promoter transaction history
	now := time.Now()
	transaction := model.PromoterTransactionHistory{
		ID:              uuid.New(),
		PromoterID:      payment.PromoterID,
		TransactionType: "EVENT_POSTING_FEE",
		Direction:       "out",
		Amount:          payment.Amount,
		ReferenceType:   stringPtr("event_creation_payment"),
		ReferenceID:     &payment.ID,
		Description:     stringPtr(fmt.Sprintf("Event creation fee for: %s", event.Title)),
		CreatedAt:       now,
	}

	if err := u.promoterTransactionRepo.Create(&transaction); err != nil {
		return errors.New("failed to create transaction record")
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}
