package repository

import (
	"context"
	"event-booking-be/internal/models"
)

type EventRepository interface {
	Create(ctx context.Context, event *models.Event) error
	GetByID(ctx context.Context, id int) (*models.Event, error)
	GetAll(ctx context.Context) ([]*models.Event, error)
	Update(ctx context.Context, id int, event *models.Event) error
	Delete(ctx context.Context, id int) error
	GetAvailableTickets(ctx context.Context, eventID int) (int, error)
	DecrementTickets(ctx context.Context, eventID int, count int) error
	IncrementTickets(ctx context.Context, eventID int, count int) error
	GetStatsByEventID(ctx context.Context, eventID int) (*models.EventStatistics, error)
	LockForUpdate(ctx context.Context, eventID int) (*models.Event, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
}

type BookingRepository interface {
	Create(ctx context.Context, booking *models.Booking) error
	GetByID(ctx context.Context, id int) (*models.Booking, error)
	GetByUserID(ctx context.Context, userID int) ([]*models.Booking, error)
	GetByEventID(ctx context.Context, eventID int) ([]*models.Booking, error)
	UpdateStatus(ctx context.Context, id int, status models.BookingStatus) error
	GetExpiredPending(ctx context.Context) ([]*models.Booking, error)
	GetWithDetails(ctx context.Context, id int) (*models.BookingWithDetails, error)
}
