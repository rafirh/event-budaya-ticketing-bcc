package usecase

import (
	"errors"
	"fmt"
	"strconv"
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
	GetByPromoterID(promoterID uuid.UUID, search, categoryID, sortBy, sortOrder string, page, limit int) ([]dto.EventPromoterListResponse, int64, error)
	CreateEvent(req dto.CreateEventRequest, promoterID uuid.UUID) (*dto.EventCreatedWithPayment, error)
	ParseCreateEventRequest(categoryIDStr, title, summaryStr, descriptionStr, venueStr, addressStr, googleMapsURLStr, latitudeStr, longitudeStr, startDateStr, endDateStr, registrationDeadlineStr, quotaStr, priceStr, isPaidStr string) (*dto.CreateEventRequest, error)
	HandleEventPaymentWebhook(req *dto.MidtransWebhookRequest) error
}

type eventUsecase struct {
	eventRepo                repository.EventRepository
	eventCreationPaymentRepo repository.EventCreationPaymentRepository
	categoryRepo             repository.EventCategoryRepository
	feeRepo                  repository.FeeRepository
	adminWalletRepo          repository.AdminWalletRepository
	promoterTransactionRepo  repository.PromoterTransactionHistoryRepository
	paymentGateway           *payment.Client
}

func NewEventUsecase(
	eventRepo repository.EventRepository,
	eventCreationPaymentRepo repository.EventCreationPaymentRepository,
	categoryRepo repository.EventCategoryRepository,
	feeRepo repository.FeeRepository,
	adminWalletRepo repository.AdminWalletRepository,
	promoterTransactionRepo repository.PromoterTransactionHistoryRepository,
	paymentGateway *payment.Client,
) EventUsecase {
	return &eventUsecase{
		eventRepo:                eventRepo,
		eventCreationPaymentRepo: eventCreationPaymentRepo,
		categoryRepo:             categoryRepo,
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

func (u *eventUsecase) GetByPromoterID(promoterID uuid.UUID, search, categoryID, sortBy, sortOrder string, page, limit int) ([]dto.EventPromoterListResponse, int64, error) {
	offset := (page - 1) * limit

	events, total, err := u.eventRepo.FindByPromoterID(promoterID, search, categoryID, sortBy, sortOrder, limit, offset)
	if err != nil {
		return nil, 0, errors.New("failed to fetch events")
	}

	responses := make([]dto.EventPromoterListResponse, 0, len(events))
	for _, event := range events {
		responses = append(responses, dto.ToEventPromoterListResponse(event, event.PaymentInfo))
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

func (u *eventUsecase) ParseCreateEventRequest(categoryIDStr, title, summaryStr, descriptionStr, venueStr, addressStr, googleMapsURLStr, latitudeStr, longitudeStr, startDateStr, endDateStr, registrationDeadlineStr, quotaStr, priceStr, isPaidStr string) (*dto.CreateEventRequest, error) {
	req := &dto.CreateEventRequest{}

	// Parse category ID
	if categoryIDStr != "" {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err != nil {
			return nil, errors.New("invalid category_id format")
		}
		req.CategoryID = &categoryID
	}

	// Parse title (required)
	if title == "" {
		return nil, errors.New("title is required")
	}
	req.Title = title

	// Parse quota (required)
	if quotaStr == "" {
		return nil, errors.New("quota is required")
	}
	quota, err := strconv.Atoi(quotaStr)
	if err != nil {
		return nil, errors.New("invalid quota format")
	}
	req.Quota = quota

	// Parse is_paid
	req.IsPaid = isPaidStr == "true"

	// Parse price (optional)
	if priceStr != "" {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return nil, errors.New("invalid price format")
		}
		req.Price = price
	}

	// Parse optional string fields
	if summaryStr != "" {
		req.Summary = &summaryStr
	}
	if descriptionStr != "" {
		req.Description = &descriptionStr
	}
	if venueStr != "" {
		req.Venue = &venueStr
	}
	if addressStr != "" {
		req.Address = &addressStr
	}
	if googleMapsURLStr != "" {
		req.GoogleMapsURL = &googleMapsURLStr
	}
	if latitudeStr != "" {
		latitude, err := strconv.ParseFloat(latitudeStr, 64)
		if err != nil {
			return nil, errors.New("invalid latitude format")
		}
		req.Latitude = &latitude
	}
	if longitudeStr != "" {
		longitude, err := strconv.ParseFloat(longitudeStr, 64)
		if err != nil {
			return nil, errors.New("invalid longitude format")
		}
		req.Longitude = &longitude
	}

	// Parse dates
	if startDateStr != "" {
		startDate, err := parseDateTime(startDateStr)
		if err != nil {
			return nil, errors.New("invalid start_date format")
		}
		req.StartDate = startDate
	}
	if endDateStr != "" {
		endDate, err := parseDateTime(endDateStr)
		if err != nil {
			return nil, errors.New("invalid end_date format")
		}
		req.EndDate = endDate
	}
	if registrationDeadlineStr != "" {
		regDeadline, err := parseDateTime(registrationDeadlineStr)
		if err != nil {
			return nil, errors.New("invalid registration_deadline format")
		}
		req.RegistrationDeadline = regDeadline
	}

	return req, nil
}

// Helper function to parse datetime strings
func parseDateTime(dateStr string) (*time.Time, error) {
	// Try ISO8601 format first
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", dateStr)
	if err == nil {
		return &t, nil
	}

	// Try simple date format
	t, err = time.Parse("2006-01-02", dateStr)
	if err == nil {
		return &t, nil
	}

	return nil, errors.New("invalid datetime format")
}

func (u *eventUsecase) CreateEvent(req dto.CreateEventRequest, promoterID uuid.UUID) (*dto.EventCreatedWithPayment, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.Quota <= 0 {
		return nil, errors.New("quota must be greater than 0")
	}

	// Validate category if provided
	if req.CategoryID != nil {
		_, err := u.categoryRepo.FindByID(*req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
	}

	fee, err := u.feeRepo.FindByType("EVENT_POSTING_FEE")
	if err != nil {
		return nil, errors.New("fee setting not found")
	}

	slug := helper.MakeSlug(req.Title)

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
		Latitude:             req.Latitude,
		Longitude:            req.Longitude,
		StartDate:            req.StartDate,
		EndDate:              req.EndDate,
		RegistrationDeadline: req.RegistrationDeadline,
		IsPaid:               req.IsPaid,
		Price:                req.Price,
		Quota:                req.Quota,
		BannerURL:            req.BannerURL,
		Status:               "draft",
		CreatedAt:            time.Now(),
	}

	if err := u.eventRepo.Create(&event); err != nil {
		return nil, errors.New("failed to create event")
	}

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
	payment, err := u.eventCreationPaymentRepo.FindByOrderID(req.OrderID)
	if err != nil {
		return errors.New("payment record not found for order ID: " + req.OrderID)
	}

	paymentStatus := mapEventPaymentStatus(req.TransactionStatus, req.FraudStatus)
	paymentMethod := req.PaymentType
	if paymentMethod != "" {
		payment.PaymentMethod = &paymentMethod
	}

	payment.Status = paymentStatus
	if paymentStatus == "paid" {
		now := time.Now()
		payment.PaidAt = &now
	}

	if err := u.eventCreationPaymentRepo.Update(payment); err != nil {
		return errors.New("failed to update payment status")
	}

	if paymentStatus != "paid" {
		return nil
	}

	event, err := u.eventRepo.FindByID(payment.EventID.String())
	if err != nil {
		return errors.New("event not found")
	}

	event.Status = "published"
	if err := u.eventRepo.Update(event); err != nil {
		return errors.New("failed to publish event")
	}

	adminWallet, err := u.adminWalletRepo.FindOrCreate()
	if err != nil {
		return errors.New("failed to get admin wallet")
	}

	adminWallet.Balance += payment.Amount
	adminWallet.TotalRevenue += payment.Amount
	if err := u.adminWalletRepo.Update(adminWallet); err != nil {
		return errors.New("failed to update admin wallet")
	}

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

func mapEventPaymentStatus(transactionStatus, fraudStatus string) string {
	switch transactionStatus {
	case "settlement":
		return "paid"
	case "capture":
		if fraudStatus == "accept" {
			return "paid"
		}
		return "pending"
	case "pending":
		return "pending"
	case "deny", "cancel", "expire", "failure":
		return "failure"
	default:
		return "pending"
	}
}

func stringPtr(s string) *string {
	return &s
}
