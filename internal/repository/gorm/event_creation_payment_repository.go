package gorm

import (
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type eventCreationPaymentRepository struct {
	db *gorm.DB
}

func NewEventCreationPaymentRepository(db *gorm.DB) repository.EventCreationPaymentRepository {
	return &eventCreationPaymentRepository{db: db}
}

func (r *eventCreationPaymentRepository) Create(payment *model.EventCreationPayment) error {
	if payment.ID == uuid.Nil {
		payment.ID = uuid.New()
	}
	return r.db.Create(payment).Error
}

func (r *eventCreationPaymentRepository) FindByOrderID(orderID string) (*model.EventCreationPayment, error) {
	var payment model.EventCreationPayment
	err := r.db.
		Preload("Event").
		Preload("Promoter").
		Where("order_id = ?", orderID).
		First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *eventCreationPaymentRepository) FindByEventID(eventID uuid.UUID) (*model.EventCreationPayment, error) {
	var payment model.EventCreationPayment
	err := r.db.
		Preload("Event").
		Preload("Promoter").
		Where("event_id = ?", eventID).
		First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *eventCreationPaymentRepository) Update(payment *model.EventCreationPayment) error {
	return r.db.Save(payment).Error
}
