package usecase

import (
	"errors"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/repository"
)

type EventUsecase interface {
	GetAll(search, categoryID, sortBy, sortOrder string, page, limit int) ([]dto.EventListResponse, int64, error)
	GetBySlug(slug string) (*dto.EventDetailResponse, error)
}

type eventUsecase struct {
	eventRepo repository.EventRepository
}

func NewEventUsecase(eventRepo repository.EventRepository) EventUsecase {
	return &eventUsecase{eventRepo: eventRepo}
}

func (u *eventUsecase) GetAll(search, categoryID, sortBy, sortOrder string, page, limit int) ([]dto.EventListResponse, int64, error) {
	offset := (page - 1) * limit

	events, total, err := u.eventRepo.FindAll(search, categoryID, sortBy, sortOrder, limit, offset)
	if err != nil {
		return nil, 0, errors.New("failed to fetch events")
	}

	responses := make([]dto.EventListResponse, 0, len(events))
	for _, event := range events {
		responses = append(responses, dto.ToEventListResponse(event))
	}

	return responses, total, nil
}

func (u *eventUsecase) GetBySlug(slug string) (*dto.EventDetailResponse, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}

	event, err := u.eventRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("event not found")
	}

	response := dto.ToEventDetailResponse(*event)
	return &response, nil
}
