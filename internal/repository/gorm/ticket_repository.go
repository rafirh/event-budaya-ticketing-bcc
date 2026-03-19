package gorm

import (
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"gorm.io/gorm"
)

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) repository.TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) CreateBatch(tickets []model.Ticket) error {
	if len(tickets) == 0 {
		return nil
	}
	return r.db.Create(&tickets).Error
}

func (r *ticketRepository) FindByOrderID(orderID string) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := r.db.Where("order_id = ?", orderID).Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}
