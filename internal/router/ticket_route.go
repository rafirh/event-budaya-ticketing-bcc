package router

import (
	"event-budaya-ticketing-bcc/internal/handler"
	"event-budaya-ticketing-bcc/internal/middleware"
	"event-budaya-ticketing-bcc/internal/repository"

	"github.com/gofiber/fiber/v2"
)

func ticketRoutes(api fiber.Router, ticketHandler *handler.TicketHandler, tokenRepo repository.PersonalAccessTokenRepository) {
	tickets := api.Group("/me/tickets", middleware.AuthMiddleware(tokenRepo))
	tickets.Get("", ticketHandler.GetMyTickets)
	tickets.Get("/:id", ticketHandler.GetMyTicketDetail)

	promoterTickets := api.Group("/events/:eventId/attendees", middleware.AuthMiddleware(tokenRepo), middleware.PromoterMiddleware())
	promoterTickets.Get("", ticketHandler.GetAttendeesByEventID)

	checkIn := api.Group("/events/:eventId/check-in", middleware.AuthMiddleware(tokenRepo), middleware.PromoterMiddleware())
	checkIn.Post("", ticketHandler.CheckInTicket)
}
