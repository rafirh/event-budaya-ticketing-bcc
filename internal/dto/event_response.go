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

type EventResponse struct {
	ID            uuid.UUID             `json:"id"`
	Promoter      EventPromoterSummary  `json:"promoter"`
	Category      *EventCategorySummary `json:"category"`
	Title         string                `json:"title"`
	Slug          *string               `json:"slug"`
	Description   *string               `json:"description"`
	Venue         *string               `json:"venue"`
	Address       *string               `json:"address"`
	GoogleMapsURL *string               `json:"google_maps_url"`
	StartDate     *time.Time            `json:"start_date"`
	EndDate       *time.Time            `json:"end_date"`
	IsPaid        bool                  `json:"is_paid"`
	BannerURL     *string               `json:"banner_url"`
	Status        string                `json:"status"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     *time.Time            `json:"updated_at"`
}

func ToEventResponse(event model.Event) EventResponse {
	var category *EventCategorySummary
	if event.Category != nil {
		category = &EventCategorySummary{
			ID:   event.Category.ID,
			Name: event.Category.Name,
		}
	}

	return EventResponse{
		ID: event.ID,
		Promoter: EventPromoterSummary{
			ID:   event.Promoter.ID,
			Name: event.Promoter.Name,
		},
		Category:      category,
		Title:         event.Title,
		Slug:          event.Slug,
		Description:   event.Description,
		Venue:         event.Venue,
		Address:       event.Address,
		GoogleMapsURL: event.GoogleMapsURL,
		StartDate:     event.StartDate,
		EndDate:       event.EndDate,
		IsPaid:        event.IsPaid,
		BannerURL:     event.BannerURL,
		Status:        event.Status,
		CreatedAt:     event.CreatedAt,
		UpdatedAt:     event.UpdatedAt,
	}
}
