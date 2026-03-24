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

func (r *ticketRepository) FindByID(id string) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.Preload("Order.Event").Where("id = ?", id).First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) FindByUserID(userID string) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := r.db.Preload("Order.Event").
		Joins("JOIN orders ON tickets.order_id = orders.id").
		Where("orders.user_id = ?", userID).
		Where("orders.status = ?", "paid").
		Order("tickets.created_at DESC").
		Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *ticketRepository) FindByEventID(eventID string, search string, limit, offset int) ([]model.Ticket, int64, error) {
	baseQuery := r.db.
		Joins("JOIN orders ON tickets.order_id = orders.id").
		Where("orders.event_id = ?", eventID).
		Where("orders.status = ?", "paid")

	if search != "" {
		searchKeyword := "%" + search + "%"
		baseQuery = baseQuery.Where(
			r.db.Where("tickets.holder_name ILIKE ?", searchKeyword).
				Or("tickets.holder_email ILIKE ?", searchKeyword).
				Or("tickets.holder_phone ILIKE ?", searchKeyword).
				Or("tickets.identity_number ILIKE ?", searchKeyword),
		)
	}

	var total int64
	if err := baseQuery.Model(&model.Ticket{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var tickets []model.Ticket
	err := baseQuery.
		Select("tickets.*").
		Order("tickets.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&tickets).Error

	if err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}
