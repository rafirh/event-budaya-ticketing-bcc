package usecase

import (
	"errors"
	"strings"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/google/uuid"
)

type EventCommentUsecase interface {
	CreateComment(userID, eventID string, req *dto.CreateEventCommentRequest) (*dto.EventCommentResponse, error)
	GetByEventID(eventID string) ([]dto.EventCommentResponse, error)
}

type eventCommentUsecase struct {
	eventRepo        repository.EventRepository
	eventCommentRepo repository.EventCommentRepository
}

func NewEventCommentUsecase(eventRepo repository.EventRepository, eventCommentRepo repository.EventCommentRepository) EventCommentUsecase {
	return &eventCommentUsecase{
		eventRepo:        eventRepo,
		eventCommentRepo: eventCommentRepo,
	}
}

func (u *eventCommentUsecase) CreateComment(userID, eventID string, req *dto.CreateEventCommentRequest) (*dto.EventCommentResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user session")
	}

	eventData, err := u.eventRepo.FindByID(eventID)
	if err != nil || eventData == nil {
		return nil, errors.New("event not found")
	}

	commentText := strings.TrimSpace(req.Comment)
	if commentText == "" {
		return nil, errors.New("comment is required")
	}

	var parentUUID *uuid.UUID
	if req.ParentID != nil && strings.TrimSpace(*req.ParentID) != "" {
		parsedParentID, err := uuid.Parse(strings.TrimSpace(*req.ParentID))
		if err != nil {
			return nil, errors.New("invalid parent_id")
		}

		parentComment, err := u.eventCommentRepo.FindByID(parsedParentID.String())
		if err != nil {
			return nil, errors.New("parent comment not found")
		}

		if parentComment.EventID != eventData.ID {
			return nil, errors.New("parent comment must belong to the same event")
		}

		if parentComment.ParentID != nil {
			return nil, errors.New("parent comment must be a top-level comment")
		}

		parentUUID = &parsedParentID
	}

	newComment := &model.EventComment{
		EventID:  eventData.ID,
		UserID:   parsedUserID,
		ParentID: parentUUID,
		Comment:  commentText,
	}

	if err := u.eventCommentRepo.Create(newComment); err != nil {
		return nil, errors.New("failed to create comment")
	}

	createdComment, err := u.eventCommentRepo.FindByID(newComment.ID.String())
	if err != nil {
		return nil, errors.New("failed to fetch created comment")
	}

	response := dto.ToEventCommentResponse(*createdComment)
	return &response, nil
}

func (u *eventCommentUsecase) GetByEventID(eventID string) ([]dto.EventCommentResponse, error) {
	eventData, err := u.eventRepo.FindByID(eventID)
	if err != nil || eventData == nil {
		return nil, errors.New("event not found")
	}

	comments, err := u.eventCommentRepo.FindByEventID(eventID)
	if err != nil {
		return nil, errors.New("failed to fetch comments")
	}

	result := make([]dto.EventCommentResponse, 0, len(comments))
	for _, comment := range comments {
		result = append(result, dto.ToEventCommentResponse(comment))
	}

	return result, nil
}
