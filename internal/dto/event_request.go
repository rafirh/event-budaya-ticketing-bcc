package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateEventRequest struct {
	CategoryID           *uuid.UUID `json:"category_id"`
	Title                string     `json:"title" validate:"required,max=200"`
	Summary              *string    `json:"summary" validate:"max=255"`
	Description          *string    `json:"description"`
	Venue                *string    `json:"venue" validate:"max=200"`
	Address              *string    `json:"address"`
	GoogleMapsURL        *string    `json:"google_maps_url"`
	StartDate            *time.Time `json:"start_date"`
	EndDate              *time.Time `json:"end_date"`
	RegistrationDeadline *time.Time `json:"registration_deadline"`
	IsPaid               bool       `json:"is_paid" validate:"required"`
	Price                float64    `json:"price" validate:"min=0"`
	Quota                int        `json:"quota" validate:"required,min=1"`
	BannerURL            *string    `json:"-"`
}

type CreateEventPaymentResponse struct {
	EventID      uuid.UUID `json:"event_id"`
	Title        string    `json:"title"`
	PaymentToken string    `json:"payment_token"`
	PaymentURL   string    `json:"payment_url"`
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
	Status       string    `json:"status"`
}

type EventCreatedWithPayment struct {
	EventID      uuid.UUID `json:"event_id"`
	Title        string    `json:"title"`
	PaymentToken string    `json:"payment_token"`
	PaymentURL   string    `json:"payment_url"`
	Amount       float64   `json:"amount"`
	Description  string    `json:"description"`
	ExpiresAt    time.Time `json:"expires_at"`
}
