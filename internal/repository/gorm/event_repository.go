package gorm

import (
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) repository.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) FindAll() ([]model.Event, error) {
	var events []model.Event
	err := r.db.
		Preload("Promoter").
		Preload("Category").
		Order("start_date ASC NULLS LAST").
		Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}
