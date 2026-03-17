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

func (r *eventCategoryRepository) FindAll(limit, offset int) ([]model.EventCategory, int64, error) {
	var total int64
	if err := r.db.Table("event_categories").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var categories []model.EventCategory
	err := r.db.
		Table("event_categories").
		Select("event_categories.id, event_categories.name, event_categories.icon, COUNT(events.id) AS event_count").
		Joins("LEFT JOIN events ON events.category_id = event_categories.id").
		Group("event_categories.id, event_categories.name, event_categories.icon").
		Order("event_categories.name ASC").
		Limit(limit).
		Offset(offset).
		Scan(&categories).Error
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}
