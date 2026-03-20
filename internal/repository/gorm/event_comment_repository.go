package gorm

import (
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"gorm.io/gorm"
)

type eventCommentRepository struct {
	db *gorm.DB
}

func NewEventCommentRepository(db *gorm.DB) repository.EventCommentRepository {
	return &eventCommentRepository{db: db}
}

func (r *eventCommentRepository) Create(comment *model.EventComment) error {
	return r.db.Create(comment).Error
}

func (r *eventCommentRepository) FindByID(id string) (*model.EventComment, error) {
	var comment model.EventComment
	err := r.db.
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Replies.User").
		Where("id = ?", id).
		First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *eventCommentRepository) FindByEventID(eventID string) ([]model.EventComment, error) {
	var comments []model.EventComment
	err := r.db.
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Replies.User").
		Where("event_id = ? AND parent_id IS NULL", eventID).
		Order("created_at DESC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}
