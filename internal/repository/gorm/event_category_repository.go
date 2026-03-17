package gorm

import (
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"gorm.io/gorm"
)

type eventCategoryRepository struct {
	db *gorm.DB
}

func NewEventCategoryRepository(db *gorm.DB) repository.EventCategoryRepository {
	return &eventCategoryRepository{db: db}
}

func (r *eventCategoryRepository) FindAll() ([]model.EventCategory, error) {
	var categories []model.EventCategory
	err := r.db.
		Table("event_categories").
		Select("event_categories.id, event_categories.name, event_categories.icon, COUNT(events.id) AS event_count").
		Joins("LEFT JOIN events ON events.category_id = event_categories.id").
		Group("event_categories.id, event_categories.name, event_categories.icon").
		Order("event_categories.name ASC").
		Scan(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
