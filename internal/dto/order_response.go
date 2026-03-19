package dto

import (
	"time"

	"github.com/google/uuid"
)

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

type MyOrderResponse struct {
	OrderID       uuid.UUID `json:"order_id"`
	EventID       uuid.UUID `json:"event_id"`
	EventName     string    `json:"event_name"`
	TicketCount   int       `json:"ticket_count"`
	TotalPrice    float64   `json:"total_price"`
	PaymentStatus string    `json:"payment_status"`
	PaymentURL    *string   `json:"payment_url"`
	OrderStatus   string    `json:"order_status"`
	CreatedAt     time.Time `json:"created_at"`
}

type TicketDetail struct {
	ID             uuid.UUID  `json:"id"`
	TicketCode     string     `json:"ticket_code"`
	HolderName     string     `json:"holder_name"`
	IdentityType   string     `json:"identity_type"`
	IdentityNumber string     `json:"identity_number"`
	HolderPhone    string     `json:"holder_phone"`
	HolderEmail    string     `json:"holder_email"`
	Notes          string     `json:"notes"`
	IsUsed         bool       `json:"is_used"`
	UsedAt         *time.Time `json:"used_at"`
}

type PaymentDetail struct {
	PaymentMethod  *string    `json:"payment_method"`
	PaymentGateway *string    `json:"payment_gateway"`
	Amount         float64    `json:"amount"`
	Status         string     `json:"status"`
	PaymentURL     *string    `json:"payment_url"`
	PaidAt         *time.Time `json:"paid_at"`
}

type MyOrderDetailResponse struct {
	OrderID     uuid.UUID      `json:"order_id"`
	EventID     uuid.UUID      `json:"event_id"`
	EventName   string         `json:"event_name"`
	TicketCount int            `json:"ticket_count"`
	UnitPrice   float64        `json:"unit_price"`
	ServiceFee  float64        `json:"service_fee"`
	TotalPrice  float64        `json:"total_price"`
	OrderStatus string         `json:"order_status"`
	CreatedAt   time.Time      `json:"created_at"`
	Tickets     []TicketDetail `json:"tickets"`
	Payment     PaymentDetail  `json:"payment"`
}

type MyTicketListResponse struct {
	ID         uuid.UUID `json:"id"`
	TicketCode string    `json:"ticket_code"`
	HolderName string    `json:"holder_name"`
	EventName  string    `json:"event_name"`
	OrderID    uuid.UUID `json:"order_id"`
	IsUsed     bool      `json:"is_used"`
	CreatedAt  time.Time `json:"created_at"`
}

type MyTicketDetailResponse struct {
	ID             uuid.UUID  `json:"id"`
	TicketCode     string     `json:"ticket_code"`
	QRCode         *string    `json:"qr_code"`
	HolderName     string     `json:"holder_name"`
	IdentityType   string     `json:"identity_type"`
	IdentityNumber string     `json:"identity_number"`
	HolderPhone    string     `json:"holder_phone"`
	HolderEmail    string     `json:"holder_email"`
	Notes          string     `json:"notes"`
	IsUsed         bool       `json:"is_used"`
	UsedAt         *time.Time `json:"used_at"`
	CreatedAt      time.Time  `json:"created_at"`
	Order          struct {
		OrderID     uuid.UUID `json:"order_id"`
		EventName   string    `json:"event_name"`
		UnitPrice   float64   `json:"unit_price"`
		OrderStatus string    `json:"order_status"`
	} `json:"order"`
}
