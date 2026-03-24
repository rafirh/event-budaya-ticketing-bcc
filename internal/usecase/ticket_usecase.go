package usecase

import (
	"errors"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/google/uuid"
)

type TicketUsecase interface {
	GetMyTickets(userID string) ([]dto.MyTicketListResponse, error)
	GetMyTicketDetail(userID, ticketID string) (*dto.MyTicketDetailResponse, error)
	GetAttendeesByEventID(eventID string, search string, page, limit int) ([]dto.AttendeeResponse, int64, error)
}

type ticketUsecase struct {
	ticketRepo repository.TicketRepository
}

func NewTicketUsecase(
	ticketRepo repository.TicketRepository,
) TicketUsecase {
	return &ticketUsecase{
		ticketRepo: ticketRepo,
	}
}

func (u *ticketUsecase) GetMyTickets(userID string) ([]dto.MyTicketListResponse, error) {
	tickets, err := u.ticketRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to fetch tickets")
	}

	result := make([]dto.MyTicketListResponse, 0, len(tickets))
	for _, ticket := range tickets {
		result = append(result, dto.MyTicketListResponse{
			ID:         ticket.ID,
			TicketCode: ticket.TicketCode,
			HolderName: ticket.HolderName,
			EventName:  ticket.Order.Event.Title,
			OrderID:    ticket.OrderID,
			IsUsed:     ticket.IsUsed,
			CreatedAt:  ticket.CreatedAt,
		})
	}

	return result, nil
}

func (u *ticketUsecase) GetMyTicketDetail(userID, ticketID string) (*dto.MyTicketDetailResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user session")
	}

	ticket, err := u.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	if ticket.Order.UserID != parsedUserID {
		return nil, errors.New("unauthorized")
	}

	if ticket.Order.Status != "paid" {
		return nil, errors.New("ticket not found")
	}

	response := &dto.MyTicketDetailResponse{
		ID:             ticket.ID,
		TicketCode:     ticket.TicketCode,
		QRCode:         ticket.QRCode,
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

	response.Order.OrderID = ticket.Order.ID
	response.Order.EventName = ticket.Order.Event.Title
	response.Order.UnitPrice = ticket.Order.UnitPrice
	response.Order.OrderStatus = ticket.Order.Status

	return response, nil
}

func (u *ticketUsecase) GetAttendeesByEventID(eventID string, search string, page, limit int) ([]dto.AttendeeResponse, int64, error) {
	offset := (page - 1) * limit

	tickets, total, err := u.ticketRepo.FindByEventID(eventID, search, limit, offset)
	if err != nil {
		return nil, 0, errors.New("failed to fetch attendees")
	}

	responses := make([]dto.AttendeeResponse, 0, len(tickets))
	for _, ticket := range tickets {
		responses = append(responses, dto.ToAttendeeResponse(ticket))
	}

	return responses, total, nil
}
