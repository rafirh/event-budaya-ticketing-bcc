package handler

import (
	"fmt"
	"html"
	"strings"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/usecase"
	"event-budaya-ticketing-bcc/pkg/response"
	"event-budaya-ticketing-bcc/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUsecase         usecase.AuthUsecase
	googleRedirectFEURI string
}

func NewAuthHandler(authUsecase usecase.AuthUsecase, googleRedirectFEURI string) *AuthHandler {
	return &AuthHandler{
		authUsecase:         authUsecase,
		googleRedirectFEURI: googleRedirectFEURI,
	}
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
	const loginURL = "https://kalcer-alpha.vercel.app/sign-in"

	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).Type("html").SendString(verificationResultHTML("Verifikasi Gagal", "Token verifikasi tidak ditemukan.", false, loginURL))
	}

	if err := h.authUsecase.VerifyEmail(token); err != nil {
		return c.Status(fiber.StatusBadRequest).Type("html").SendString(verificationResultHTML("Verifikasi Gagal", err.Error(), false, loginURL))
	}

	return c.Status(fiber.StatusOK).Type("html").SendString(verificationResultHTML("Verifikasi Berhasil", "Akun Anda berhasil diverifikasi.", true, loginURL))
}

func verificationResultHTML(title, message string, success bool, loginURL string) string {
	statusClass := "status-fail"
	if success {
		statusClass = "status-ok"
	}

	return fmt.Sprintf(`<!doctype html>
<html lang="id">
<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<title>%s</title>
	<style>
		:root {
			--bg: #f4f7fb;
			--card: #ffffff;
			--text: #17212f;
			--muted: #5a6678;
			--ok: #16803a;
			--fail: #b42318;
			--btn: #9f7f47;
			--btn-hover: #8b6f3e;
		}
		* { box-sizing: border-box; }
		body {
			margin: 0;
			min-height: 100vh;
			display: grid;
			place-items: center;
			background: radial-gradient(circle at top, #eaf2ff 0%%, var(--bg) 45%%);
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
			color: var(--text);
			padding: 20px;
		}
		.card {
			width: 100%%;
			max-width: 520px;
			background: var(--card);
			border-radius: 16px;
			padding: 28px;
			box-shadow: 0 10px 30px rgba(16, 24, 40, 0.08);
			border: 1px solid #e5e7eb;
			text-align: center;
		}
		h1 {
			margin: 0 0 10px;
			font-size: 28px;
			line-height: 1.2;
		}
		p {
			margin: 0;
			color: var(--muted);
			line-height: 1.6;
			font-size: 15px;
		}
		.status-ok { color: var(--ok); }
		.status-fail { color: var(--fail); }
		.btn {
			margin-top: 24px;
			display: inline-block;
			background: var(--btn);
			color: #fff;
			text-decoration: none;
			font-weight: 600;
			padding: 12px 20px;
			border-radius: 10px;
			transition: background 120ms ease-in-out;
		}
		.btn:hover { background: var(--btn-hover); }
	</style>
</head>
<body>
	<main class="card">
		<h1 class="%s">%s</h1>
		<p>%s</p>
		<a class="btn" href="%s">Ke Halaman Login</a>
	</main>
</body>
</html>`,
		html.EscapeString(title),
		statusClass,
		html.EscapeString(title),
		html.EscapeString(message),
		html.EscapeString(loginURL),
	)
}

func (h *AuthHandler) ResendVerificationEmail(c *fiber.Ctx) error {
	var req dto.ResendVerificationEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return response.ValidationError(c, errors)
	}

	if err := h.authUsecase.ResendVerificationEmail(req.Email); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Verification email resent successfully", nil)
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

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return response.ValidationError(c, errors)
	}

	userID := c.Locals("userID").(string)
	if err := h.authUsecase.ChangePassword(userID, &req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Password changed successfully", nil)
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

func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	state, err := dto.GenerateRandomState()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to generate state")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   600,
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	url := h.authUsecase.GoogleLoginURL(state)
	return c.Redirect(url)
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return response.Error(c, fiber.StatusBadRequest, "Missing code or state parameter")
	}

	// Verify state from cookie
	storedState := c.Cookies("oauth_state")
	if storedState != state {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid state parameter - CSRF protection failed")
	}

	// Exchange code for token and login
	result, err := h.authUsecase.GoogleCallback(c.Context(), code, state)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	// Redirect to FE URL with token as query parameter
	redirectURL := fmt.Sprintf("%s?token=%s", h.googleRedirectFEURI, result.Token)
	return c.Redirect(redirectURL)
}
