package dto

type CheckInRequest struct {
	TicketCode string `json:"ticket_code" validate:"required"`
}

type CheckInResponse struct {
	TicketCode string  `json:"ticket_code"`
	HolderName string  `json:"holder_name"`
	IsUsed     bool    `json:"is_used"`
	UsedAt     *string `json:"used_at,omitempty"`
}
