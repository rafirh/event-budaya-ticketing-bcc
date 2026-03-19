package dto

type TicketHolderRequest struct {
	HolderName     string  `json:"holder_name" validate:"required"`
	IdentityType   string  `json:"identity_type" validate:"required"`
	IdentityNumber string  `json:"identity_number" validate:"required"`
	HolderPhone    string  `json:"holder_phone" validate:"required"`
	HolderEmail    string  `json:"holder_email" validate:"required,email"`
	Notes          *string `json:"notes"`
}

type CreateTicketOrderRequest struct {
	EventID string                `json:"event_id" validate:"required,uuid"`
	Tickets []TicketHolderRequest `json:"tickets" validate:"required,min=1,dive"`
}

type MidtransWebhookRequest struct {
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
	TransactionID     string `json:"transaction_id"`
}
