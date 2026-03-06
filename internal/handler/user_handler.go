package handler

import (
	"strconv"

	"event-budaya-ticketing-bcc/internal/domain"
	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"
	"event-budaya-ticketing-bcc/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

// Register handles user registration
// @Summary Register new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.CreateUserRequest true "User registration data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req domain.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if errors := validator.ValidateStruct(req); errors != nil {
		return response.ValidationError(c, errors)
	}

	user, err := h.userUsecase.Register(&req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "User registered successfully", user)
}

// Login handles user authentication
// @Summary User login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "User login credentials"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if errors := validator.ValidateStruct(req); errors != nil {
		return response.ValidationError(c, errors)
	}

	result, err := h.userUsecase.Login(&req)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Login successful", result)
}

// GetAllUsers returns all users
// @Summary Get all users
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userUsecase.GetAllUsers()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Users retrieved successfully", users)
}

// GetUserByID returns a user by ID
// @Summary Get user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.userUsecase.GetUserByID(uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "User retrieved successfully", user)
}

// UpdateUser updates a user by ID
// @Summary Update user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body domain.UpdateUserRequest true "User update data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var req domain.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if errors := validator.ValidateStruct(req); errors != nil {
		return response.ValidationError(c, errors)
	}

	user, err := h.userUsecase.UpdateUser(uint(id), &req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "User updated successfully", user)
}

// DeleteUser deletes a user by ID
// @Summary Delete user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := h.userUsecase.DeleteUser(uint(id)); err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "User deleted successfully", nil)
}

// GetProfile returns the authenticated user's profile
// @Summary Get current user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/auth/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	user, err := h.userUsecase.GetUserByID(userID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Profile retrieved successfully", user)
}
