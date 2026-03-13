package usecase

import (
	"errors"

	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"
)

type CategoryUsecase interface {
	GetAll() ([]model.EventCategory, error)
}

type categoryUsecase struct {
	categoryRepo repository.EventCategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.EventCategoryRepository) CategoryUsecase {
	return &categoryUsecase{categoryRepo: categoryRepo}
}

func (u *categoryUsecase) GetAll() ([]model.EventCategory, error) {
	categories, err := u.categoryRepo.FindAll()
	if err != nil {
		return nil, errors.New("failed to fetch categories")
	}
	return categories, nil
}
