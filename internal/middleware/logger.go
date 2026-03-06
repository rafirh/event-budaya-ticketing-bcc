package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggerMiddleware logs incoming requests
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log request details
		log.Printf(
			"[%s] %s %s - %d - %v",
			c.Method(),
			c.Path(),
			c.IP(),
			c.Response().StatusCode(),
			time.Since(start),
		)

		return err
	}
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v", r)
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "Internal server error",
				})
			}
		}()
		return c.Next()
	}
}
