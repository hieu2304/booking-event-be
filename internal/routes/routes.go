package routes

import (
	"event-booking-be/internal/handler"
	"event-booking-be/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	userHandler    *handler.UserHandler
	eventHandler   *handler.EventHandler
	bookingHandler *handler.BookingHandler
}

func NewRouter(
	userHandler *handler.UserHandler,
	eventHandler *handler.EventHandler,
	bookingHandler *handler.BookingHandler,
) *Router {
	return &Router{
		userHandler:    userHandler,
		eventHandler:   eventHandler,
		bookingHandler: bookingHandler,
	}
}

func (r *Router) Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/register", r.userHandler.Register)
	auth.Post("/login", r.userHandler.Login)

	// Event routes (public read, protected write)
	events := api.Group("/events")
	events.Get("/", r.eventHandler.GetAllEvents)
	events.Get("/:id", r.eventHandler.GetEvent)
	events.Get("/:id/statistics", r.eventHandler.GetEventStatistics)
	events.Post("/", r.eventHandler.CreateEvent)
	events.Put("/:id", r.eventHandler.UpdateEvent)
	events.Delete("/:id", r.eventHandler.DeleteEvent)

	// Protected booking routes
	bookings := api.Group("/bookings", middleware.AuthMiddleware())
	bookings.Post("/", r.bookingHandler.CreateBooking)
	bookings.Get("/", r.bookingHandler.GetUserBookings)
	bookings.Get("/:id", r.bookingHandler.GetBooking)
	bookings.Post("/:id/confirm", r.bookingHandler.ConfirmPayment)
	bookings.Post("/:id/cancel", r.bookingHandler.CancelBooking)

	// Protected user routes
	users := api.Group("/users", middleware.AuthMiddleware())
	users.Get("/profile", r.userHandler.GetProfile)
}
