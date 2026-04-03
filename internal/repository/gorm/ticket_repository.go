package gorm

import (
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"
	"time"

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

func (r *ticketRepository) FindDistinctEventsByUserIDAndDateRange(userID string, start, end time.Time) ([]model.Event, error) {
	var events []model.Event
	err := r.db.Model(&model.Event{}).
		Distinct("events.*").
		Joins("JOIN orders ON orders.event_id = events.id").
		Joins("JOIN tickets ON tickets.order_id = orders.id").
		Where("orders.user_id = ?", userID).
		Where("orders.status = ?", "paid").
		Where("events.start_date IS NOT NULL").
		Where("events.start_date >= ?", start).
		Where("events.start_date < ?", end).
		Order("events.start_date ASC").
		Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
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

func (r *ticketRepository) FindByIDAndEventID(ticketID, eventID string) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.
		Joins("JOIN orders ON tickets.order_id = orders.id").
		Where("tickets.id = ?", ticketID).
		Where("orders.event_id = ?", eventID).
		First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) FindByCodeAndEventID(ticketCode, eventID string) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.
		Joins("JOIN orders ON tickets.order_id = orders.id").
		Where("tickets.ticket_code = ?", ticketCode).
		Where("orders.event_id = ?", eventID).
		First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) Update(ticket *model.Ticket) error {
	return r.db.Save(ticket).Error
}
