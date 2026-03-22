package handler

import (
	"strings"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"
	"event-budaya-ticketing-bcc/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return response.ValidationError(c, errors)
	}

	user, err := h.authUsecase.Register(&req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "User registered successfully, please check your email for account activation", user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return response.ValidationError(c, errors)
	}

	result, err := h.authUsecase.Login(&req)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Login successful", result)
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return response.Error(c, fiber.StatusBadRequest, "verification token is required")
	}

	if err := h.authUsecase.VerifyEmail(token); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Email verified successfully. You can now login.", nil)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	user, err := h.authUsecase.GetMe(userID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "User profile retrieved", user)
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	if c.FormValue("email") != "" {
		return response.Error(c, fiber.StatusBadRequest, "email cannot be updated")
	}

	var req dto.UpdateProfileRequest

	if name := c.FormValue("name"); name != "" {
		req.Name = &name
	}
	if phone := c.FormValue("phone"); phone != "" {
		req.Phone = &phone
	}
	if gender := c.FormValue("gender"); gender != "" {
		req.Gender = &gender
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return response.ValidationError(c, errors)
	}

	fileHeader, _ := c.FormFile("profile_photo")
	userID := c.Locals("userID").(string)

	user, err := h.authUsecase.UpdateProfile(userID, &req, fileHeader)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Profile updated successfully", user)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return response.Error(c, fiber.StatusBadRequest, "Invalid authorization header")
	}

	if err := h.authUsecase.Logout(parts[1]); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Logged out successfully", nil)
}
