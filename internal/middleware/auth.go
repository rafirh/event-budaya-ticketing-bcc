package middleware

import (
	"strings"

	"event-budaya-ticketing-bcc/config"
	"event-budaya-ticketing-bcc/pkg/jwt"
	"event-budaya-ticketing-bcc/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT token
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Missing authorization header")
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		tokenString := parts[1]

		// Validate token
		claims, err := jwt.ValidateToken(tokenString, config.AppConfig.JWTSecret)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// AdminMiddleware checks if user has admin role
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("userRole")
		if role == nil || role.(string) != "admin" {
			return response.Error(c, fiber.StatusForbidden, "Access denied: admin role required")
		}
		return c.Next()
	}
}

// OptionalAuthMiddleware validates JWT token if present but doesn't fail if missing
func OptionalAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := parts[1]

		claims, err := jwt.ValidateToken(tokenString, config.AppConfig.JWTSecret)
		if err == nil {
			c.Locals("userID", claims.UserID)
			c.Locals("userEmail", claims.Email)
			c.Locals("userRole", claims.Role)
		}

		return c.Next()
	}
}
