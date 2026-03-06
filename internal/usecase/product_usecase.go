package usecase

import (
	"errors"

	"event-budaya-ticketing-bcc/internal/domain"
	"event-budaya-ticketing-bcc/internal/repository"
)

type ProductUsecase interface {
	CreateProduct(product *domain.Product) (*domain.Product, error)
	GetAllProducts() ([]domain.Product, error)
	GetProductByID(id uint) (*domain.Product, error)
	UpdateProduct(id uint, product *domain.Product) (*domain.Product, error)
	DeleteProduct(id uint) error
}

type productUsecase struct {
	productRepo repository.ProductRepository
}

// NewProductUsecase creates a new instance of ProductUsecase
func NewProductUsecase(productRepo repository.ProductRepository) ProductUsecase {
	return &productUsecase{productRepo: productRepo}
}

func (u *productUsecase) CreateProduct(product *domain.Product) (*domain.Product, error) {
	if err := u.productRepo.Create(product); err != nil {
		return nil, errors.New("failed to create product")
	}
	return product, nil
}

func (u *productUsecase) GetAllProducts() ([]domain.Product, error) {
	products, err := u.productRepo.FindAll()
	if err != nil {
		return nil, errors.New("failed to fetch products")
	}
	return products, nil
}

func (u *productUsecase) GetProductByID(id uint) (*domain.Product, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (u *productUsecase) UpdateProduct(id uint, req *domain.Product) (*domain.Product, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	if err := u.productRepo.Update(product); err != nil {
		return nil, errors.New("failed to update product")
	}

	return product, nil
}

func (u *productUsecase) DeleteProduct(id uint) error {
	_, err := u.productRepo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	if err := u.productRepo.Delete(id); err != nil {
		return errors.New("failed to delete product")
	}

	return nil
}
