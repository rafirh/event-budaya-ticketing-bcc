package usecase

import (
	"errors"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/repository"
)

type EventUsecase interface {
	GetAll() ([]dto.EventResponse, error)
	GetBySlug(slug string) (*dto.EventResponse, error)
}

type eventUsecase struct {
	eventRepo repository.EventRepository
}

func NewEventUsecase(eventRepo repository.EventRepository) EventUsecase {
	return &eventUsecase{eventRepo: eventRepo}
}

func (u *eventUsecase) GetAll() ([]dto.EventResponse, error) {
	events, err := u.eventRepo.FindAll()
	if err != nil {
		return nil, errors.New("failed to fetch events")
	}

	responses := make([]dto.EventResponse, 0, len(events))
	for _, event := range events {
		responses = append(responses, dto.ToEventResponse(event))
	}

	return responses, nil
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
