package usecase

import (
	"errors"

	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"
)

type CategoryUsecase interface {
	GetAll(page, limit int, sortBy, sortOrder string) ([]model.EventCategory, int64, error)
}

type categoryUsecase struct {
	categoryRepo repository.EventCategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.EventCategoryRepository) CategoryUsecase {
	return &categoryUsecase{categoryRepo: categoryRepo}
}

func (u *categoryUsecase) GetAll(page, limit int, sortBy, sortOrder string) ([]model.EventCategory, int64, error) {
	offset := (page - 1) * limit

	categories, total, err := u.categoryRepo.FindAll(limit, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, errors.New("failed to fetch categories")
	}
	return categories, total, nil
}
