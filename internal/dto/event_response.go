package dto

import (
	"time"

	"event-budaya-ticketing-bcc/internal/model"

	"github.com/google/uuid"
)

type EventCategorySummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type EventPromoterSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type EventListResponse struct {
	ID                   uuid.UUID             `json:"id"`
	Category             *EventCategorySummary `json:"category"`
	BannerURL            *string               `json:"banner_url"`
	Slug                 *string               `json:"slug"`
	Summary              *string               `json:"summary"`
	StartDate            *time.Time            `json:"start_date"`
	EndDate              *time.Time            `json:"end_date"`
	RegistrationDeadline *time.Time            `json:"registration_deadline"`
	Venue                *string               `json:"venue"`
	Title                string                `json:"title"`
	Price                float64               `json:"price"`
	Quota                int                   `json:"quota"`
	Sold                 int                   `json:"sold"`
}

type EventDetailResponse struct {
	ID                   uuid.UUID             `json:"id"`
	Promoter             EventPromoterSummary  `json:"promoter"`
	Category             *EventCategorySummary `json:"category"`
	Title                string                `json:"title"`
	Slug                 *string               `json:"slug"`
	Summary              *string               `json:"summary"`
	Description          *string               `json:"description"`
	Venue                *string               `json:"venue"`
	Address              *string               `json:"address"`
	GoogleMapsURL        *string               `json:"google_maps_url"`
	Latitude             *float64              `json:"latitude"`
	Longitude            *float64              `json:"longitude"`
	StartDate            *time.Time            `json:"start_date"`
	EndDate              *time.Time            `json:"end_date"`
	RegistrationDeadline *time.Time            `json:"registration_deadline"`
	IsPaid               bool                  `json:"is_paid"`
	Price                float64               `json:"price"`
	Quota                int                   `json:"quota"`
	Sold                 int                   `json:"sold"`
	BannerURL            *string               `json:"banner_url"`
}

func ToEventListResponse(event model.Event) EventListResponse {
	var category *EventCategorySummary
	if event.Category != nil {
		category = &EventCategorySummary{
			ID:   event.Category.ID,
			Name: event.Category.Name,
		}
	}

	return EventListResponse{
		ID:                   event.ID,
		Category:             category,
		BannerURL:            event.BannerURL,
		Slug:                 event.Slug,
		Title:                event.Title,
		Summary:              event.Summary,
		StartDate:            event.StartDate,
		EndDate:              event.EndDate,
		RegistrationDeadline: event.RegistrationDeadline,
		Venue:                event.Venue,
		Price:                event.Price,
		Quota:                event.Quota,
		Sold:                 event.Sold,
	}
}

func ToEventDetailResponse(event model.Event) EventDetailResponse {
	var category *EventCategorySummary
	if event.Category != nil {
		category = &EventCategorySummary{
			ID:   event.Category.ID,
			Name: event.Category.Name,
		}
	}

	return EventDetailResponse{
		ID: event.ID,
		Promoter: EventPromoterSummary{
			ID:   event.Promoter.ID,
			Name: event.Promoter.Name,
		},
		Category:             category,
		Title:                event.Title,
		Slug:                 event.Slug,
		Summary:              event.Summary,
		Description:          event.Description,
		Venue:                event.Venue,
		Address:              event.Address,
		GoogleMapsURL:        event.GoogleMapsURL,
		Latitude:             event.Latitude,
		Longitude:            event.Longitude,
		StartDate:            event.StartDate,
		EndDate:              event.EndDate,
		RegistrationDeadline: event.RegistrationDeadline,
		IsPaid:               event.IsPaid,
		Price:                event.Price,
		Quota:                event.Quota,
		Sold:                 event.Sold,
		BannerURL:            event.BannerURL,
	}
}

type PaymentInfo struct {
	Status        string     `json:"status"` // pending, settlement, failure, cancelled
	Amount        float64    `json:"amount"`
	PaymentURL    *string    `json:"payment_url"`
	PaymentMethod *string    `json:"payment_method"`
	PaidAt        *time.Time `json:"paid_at"`
}

type EventPromoterListResponse struct {
	ID                   uuid.UUID             `json:"id"`
	Title                string                `json:"title"`
	Slug                 *string               `json:"slug"`
	Category             *EventCategorySummary `json:"category"`
	BannerURL            *string               `json:"banner_url"`
	Summary              *string               `json:"summary"`
	Status               string                `json:"status"` // draft, published, cancelled
	StartDate            *time.Time            `json:"start_date"`
	EndDate              *time.Time            `json:"end_date"`
	RegistrationDeadline *time.Time            `json:"registration_deadline"`
	Venue                *string               `json:"venue"`
	IsPaid               bool                  `json:"is_paid"`
	Price                float64               `json:"price"`
	Quota                int                   `json:"quota"`
	Sold                 int                   `json:"sold"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            *time.Time            `json:"updated_at"`
	Payment              *PaymentInfo          `json:"payment"`
}

func ToEventPromoterListResponse(event model.Event, payment *model.EventCreationPayment) EventPromoterListResponse {
	var category *EventCategorySummary
	if event.Category != nil {
		category = &EventCategorySummary{
			ID:   event.Category.ID,
			Name: event.Category.Name,
		}
	}

	var paymentInfo *PaymentInfo
	if payment != nil {
		paymentInfo = &PaymentInfo{
			Status:        payment.Status,
			Amount:        payment.Amount,
			PaymentURL:    payment.PaymentURL,
			PaymentMethod: payment.PaymentMethod,
			PaidAt:        payment.PaidAt,
		}
	}

	return EventPromoterListResponse{
		ID:                   event.ID,
		Title:                event.Title,
		Slug:                 event.Slug,
		Category:             category,
		BannerURL:            event.BannerURL,
		Summary:              event.Summary,
		Status:               event.Status,
		StartDate:            event.StartDate,
		EndDate:              event.EndDate,
		RegistrationDeadline: event.RegistrationDeadline,
		Venue:                event.Venue,
		IsPaid:               event.IsPaid,
		Price:                event.Price,
		Quota:                event.Quota,
		Sold:                 event.Sold,
		CreatedAt:            event.CreatedAt,
		UpdatedAt:            event.UpdatedAt,
		Payment:              paymentInfo,
	}
}
