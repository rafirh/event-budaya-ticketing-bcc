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
	Category             *EventCategorySummary `json:"category"`
	BannerURL            *string               `json:"banner_url"`
	Slug                 *string               `json:"slug"`
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
	Description          *string               `json:"description"`
	Venue                *string               `json:"venue"`
	Address              *string               `json:"address"`
	GoogleMapsURL        *string               `json:"google_maps_url"`
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
		Category:             category,
		BannerURL:            event.BannerURL,
		Slug:                 event.Slug,
		StartDate:            event.StartDate,
		EndDate:              event.EndDate,
		RegistrationDeadline: event.RegistrationDeadline,
		Venue:                event.Venue,
		Title:                event.Title,
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
		Description:          event.Description,
		Venue:                event.Venue,
		Address:              event.Address,
		GoogleMapsURL:        event.GoogleMapsURL,
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
