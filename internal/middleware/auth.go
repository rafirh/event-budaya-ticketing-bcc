package middleware

import (
	"strings"
	"time"

	"event-budaya-ticketing-bcc/internal/repository"
	"event-budaya-ticketing-bcc/pkg/response"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(tokenRepo repository.PersonalAccessTokenRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		tokenString := parts[1]

		pat, err := tokenRepo.FindByToken(tokenString)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token")
			}
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		now := time.Now()
		pat.LastUsedAt = &now
		_ = tokenRepo.UpdateLastUsed(pat.ID, now)

		c.Locals("userID", pat.User.ID.String())
		c.Locals("userEmail", pat.User.Email)
		c.Locals("userRole", pat.User.Role)

		return c.Next()
	}
}

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("userRole")
		if role == nil || role.(string) != "admin" {
			return response.Error(c, fiber.StatusForbidden, "Access denied: admin role required")
		}
		return c.Next()
	}
}

func PromoterMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("userRole")
		if role == nil || role.(string) != "promotor" {
			return response.Error(c, fiber.StatusForbidden, "Access denied: promoter role required")
		}
		return c.Next()
	}
}
