package dto

import (
	"time"

	"event-budaya-ticketing-bcc/internal/model"
)

type AttendeeResponse struct {
	ID             string     `json:"id"`
	TicketCode     string     `json:"ticket_code"`
	HolderName     string     `json:"holder_name"`
	IdentityType   string     `json:"identity_type"`
	IdentityNumber string     `json:"identity_number"`
	HolderPhone    string     `json:"holder_phone"`
	HolderEmail    string     `json:"holder_email"`
	Notes          string     `json:"notes"`
	IsUsed         bool       `json:"is_used"`
	UsedAt         *time.Time `json:"used_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

func ToAttendeeResponse(ticket model.Ticket) AttendeeResponse {
	return AttendeeResponse{
		ID:             ticket.ID.String(),
		TicketCode:     ticket.TicketCode,
		HolderName:     ticket.HolderName,
		IdentityType:   ticket.IdentityType,
		IdentityNumber: ticket.IdentityNumber,
		HolderPhone:    ticket.HolderPhone,
		HolderEmail:    ticket.HolderEmail,
		Notes:          ticket.Notes,
		IsUsed:         ticket.IsUsed,
		UsedAt:         ticket.UsedAt,
		CreatedAt:      ticket.CreatedAt,
	}
}
