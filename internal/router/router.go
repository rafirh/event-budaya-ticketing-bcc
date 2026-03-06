package router

import (
	"event-budaya-ticketing-bcc/internal/handler"
	"event-budaya-ticketing-bcc/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler, productHandler *handler.ProductHandler) {
	// Global middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.RecoveryMiddleware())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API routes
	api := app.Group("/api")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	// Protected auth routes
	authProtected := auth.Group("", middleware.AuthMiddleware())
	authProtected.Get("/profile", userHandler.GetProfile)

	// User routes (protected)
	users := api.Group("/users", middleware.AuthMiddleware())
	users.Get("/", userHandler.GetAllUsers)
	users.Get("/:id", userHandler.GetUserByID)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", middleware.AdminMiddleware(), userHandler.DeleteUser)

	// Product routes
	products := api.Group("/products")
	products.Get("/", productHandler.GetAllProducts)
	products.Get("/:id", productHandler.GetProductByID)

	// Protected product routes (admin only)
	productsProtected := products.Group("", middleware.AuthMiddleware(), middleware.AdminMiddleware())
	productsProtected.Post("/", productHandler.CreateProduct)
	productsProtected.Put("/:id", productHandler.UpdateProduct)
	productsProtected.Delete("/:id", productHandler.DeleteProduct)
}
