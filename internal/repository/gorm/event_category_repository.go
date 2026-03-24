package gorm

import (
	"fmt"

	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type eventCategoryRepository struct {
	db *gorm.DB
}

func NewEventCategoryRepository(db *gorm.DB) repository.EventCategoryRepository {
	return &eventCategoryRepository{db: db}
}

func (r *eventCategoryRepository) FindAll(limit, offset int, sortBy, sortOrder string) ([]model.EventCategory, int64, error) {
	var total int64
	if err := r.db.Table("event_categories").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderBy := buildCategoryOrderBy(sortBy, sortOrder)

	var categories []model.EventCategory
	err := r.db.
		Table("event_categories").
		Select("event_categories.id, event_categories.name, event_categories.icon, COUNT(events.id) AS event_count").
		Joins("LEFT JOIN events ON events.category_id = event_categories.id").
		Group("event_categories.id, event_categories.name, event_categories.icon").
		Order(orderBy).
		Limit(limit).
		Offset(offset).
		Scan(&categories).Error
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (r *eventCategoryRepository) FindByID(id uuid.UUID) (*model.EventCategory, error) {
	var category model.EventCategory
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func buildCategoryOrderBy(sortBy, sortOrder string) string {
	column := map[string]string{
		"name":        "event_categories.name",
		"event_count": "event_count",
	}[sortBy]

	if column == "" {
		column = "event_categories.name"
	}

	order := "ASC"
	if sortOrder == "desc" {
		order = "DESC"
	}

	return fmt.Sprintf("%s %s", column, order)
}
