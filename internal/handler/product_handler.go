package handler

import (
	"strconv"

	"event-budaya-ticketing-bcc/internal/domain"
	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

// NewProductHandler creates a new instance of ProductHandler
func NewProductHandler(productUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: productUsecase}
}

// CreateProduct creates a new product
// @Summary Create product
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.Product true "Product data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var product domain.Product
	if err := c.BodyParser(&product); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.productUsecase.CreateProduct(&product)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Product created successfully", result)
}

// GetAllProducts returns all products
// @Summary Get all products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/products [get]
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productUsecase.GetAllProducts()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Products retrieved successfully", products)
}

// GetProductByID returns a product by ID
// @Summary Get product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	product, err := h.productUsecase.GetProductByID(uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Product retrieved successfully", product)
}

// UpdateProduct updates a product by ID
// @Summary Update product
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body domain.Product true "Product data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	var product domain.Product
	if err := c.BodyParser(&product); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.productUsecase.UpdateProduct(uint(id), &product)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Product updated successfully", result)
}

// DeleteProduct deletes a product by ID
// @Summary Delete product
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	if err := h.productUsecase.DeleteProduct(uint(id)); err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Product deleted successfully", nil)
}
