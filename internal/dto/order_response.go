package dto

import "github.com/google/uuid"

type CreateTicketOrderResponse struct {
	OrderID         uuid.UUID `json:"order_id"`
	EventID         uuid.UUID `json:"event_id"`
	TicketCount     int       `json:"ticket_count"`
	UnitPrice       float64   `json:"unit_price"`
	ServiceFee      float64   `json:"service_fee"`
	ServiceFeeTotal float64   `json:"service_fee_total"`
	TotalPrice      float64   `json:"total_price"`
	PaymentStatus   string    `json:"payment_status"`
	PaymentToken    string    `json:"payment_token"`
	PaymentURL      string    `json:"payment_url"`
}
