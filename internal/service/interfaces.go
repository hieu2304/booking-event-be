package service

import (
	"context"
	"event-booking-be/internal/models"
)

type EventService interface {
	CreateEvent(ctx context.Context, req *models.CreateEventRequest) (*models.Event, error)
	GetEvent(ctx context.Context, id int) (*models.Event, error)
	GetAllEvents(ctx context.Context) ([]*models.Event, error)
	UpdateEvent(ctx context.Context, id int, req *models.UpdateEventRequest) (*models.Event, error)
	DeleteEvent(ctx context.Context, id int) error
	GetEventStatistics(ctx context.Context, eventID int) (*models.EventStatistics, error)
}

type BookingService interface {
	CreateBooking(ctx context.Context, userID int, req *models.CreateBookingRequest) (*models.Booking, error)
	GetBooking(ctx context.Context, id int) (*models.BookingWithDetails, error)
	GetUserBookings(ctx context.Context, userID int) ([]*models.Booking, error)
	ConfirmPayment(ctx context.Context, bookingID int) error
	CancelBooking(ctx context.Context, bookingID int) error
	ProcessExpiredBookings(ctx context.Context) error
}

type UserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.User, error)
}
