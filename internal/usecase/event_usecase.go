package usecase

import (
	"errors"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/repository"
)

type EventUsecase interface {
	GetAll(search, categoryID string, page, limit int) ([]dto.EventResponse, int64, error)
	GetBySlug(slug string) (*dto.EventResponse, error)
}

type eventUsecase struct {
	eventRepo repository.EventRepository
}

func NewEventUsecase(eventRepo repository.EventRepository) EventUsecase {
	return &eventUsecase{eventRepo: eventRepo}
}

func (u *eventUsecase) GetAll(search, categoryID string, page, limit int) ([]dto.EventResponse, int64, error) {
	offset := (page - 1) * limit

	events, total, err := u.eventRepo.FindAll(search, categoryID, limit, offset)
	if err != nil {
		return nil, 0, errors.New("failed to fetch events")
	}

	responses := make([]dto.EventResponse, 0, len(events))
	for _, event := range events {
		responses = append(responses, dto.ToEventResponse(event))
	}

	return responses, total, nil
}

func (u *eventUsecase) GetBySlug(slug string) (*dto.EventResponse, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}

	event, err := u.eventRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("event not found")
	}

	response := dto.ToEventResponse(*event)
	return &response, nil
}
