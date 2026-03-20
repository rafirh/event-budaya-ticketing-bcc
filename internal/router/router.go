package router

import (
	"event-budaya-ticketing-bcc/internal/handler"
	"event-budaya-ticketing-bcc/internal/middleware"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App, authHandler *handler.AuthHandler, categoryHandler *handler.CategoryHandler, eventHandler *handler.EventHandler, orderHandler *handler.OrderHandler, ticketHandler *handler.TicketHandler, eventCommentHandler *handler.EventCommentHandler, walletHandler *handler.WalletHandler, tokenRepo repository.PersonalAccessTokenRepository) {
	app.Static("/uploads", "./public/uploads")
	app.Static("/", "./public")

	registerMiddlewares(app)
	registerBaseRoutes(app)

	api := app.Group("/api")
	authRoutes(api, authHandler, tokenRepo)
	categoryRoutes(api, categoryHandler)
	eventRoutes(api, eventHandler)
	commentRoutes(api, eventCommentHandler, tokenRepo)
	orderRoutes(api, orderHandler, tokenRepo)
	ticketRoutes(api, ticketHandler, tokenRepo)
	SetupWalletRoute(api, walletHandler, tokenRepo)
	webhookRoutes(api, orderHandler)
}

func registerMiddlewares(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.RecoveryMiddleware())
}

func registerBaseRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":  "API running successfully",
			"status":   true,
			"api_docs": "https://2ikhh28mj3.apidog.io/",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})
}
