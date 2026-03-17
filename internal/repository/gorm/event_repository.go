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

func (r *eventRepository) FindAll(search, categoryID string, limit, offset int) ([]model.Event, int64, error) {
	baseQuery := r.db.Model(&model.Event{})

	if search != "" {
		searchKeyword := "%" + search + "%"
		baseQuery = baseQuery.Where("title ILIKE ? OR summary ILIKE ? OR description ILIKE ?", searchKeyword, searchKeyword, searchKeyword)
	}

	if categoryID != "" {
		baseQuery = baseQuery.Where("category_id = ?", categoryID)
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var events []model.Event
	err := r.db.
		Preload("Promoter").
		Preload("Category").
		Scopes(func(db *gorm.DB) *gorm.DB {
			query := db
			if search != "" {
				searchKeyword := "%" + search + "%"
				query = query.Where("title ILIKE ? OR summary ILIKE ? OR description ILIKE ?", searchKeyword, searchKeyword, searchKeyword)
			}
			if categoryID != "" {
				query = query.Where("category_id = ?", categoryID)
			}
			return query
		}).
		Order("start_date ASC NULLS LAST").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	if err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

func (r *eventRepository) FindBySlug(slug string) (*model.Event, error) {
	var event model.Event
	err := r.db.
		Preload("Promoter").
		Preload("Category").
		Where("slug = ?", slug).
		First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}
