package service

import (
	"context"
	"event-booking-be/internal/models"
	"event-booking-be/internal/repository"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type bookingService struct {
	bookingRepo repository.BookingRepository
	eventRepo   repository.EventRepository
	db          *gorm.DB
	timeout     time.Duration
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	eventRepo repository.EventRepository,
	db *gorm.DB,
	timeoutMinutes int,
) BookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
		db:          db,
		timeout:     time.Duration(timeoutMinutes) * time.Minute,
	}
}

func (s *bookingService) CreateBooking(ctx context.Context, userID int, req *models.CreateBookingRequest) (*models.Booking, error) {
	var booking *models.Booking
	
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock event row
		event, err := s.eventRepo.LockForUpdate(ctx, req.EventID)
		if err != nil {
			return fmt.Errorf("event not found")
		}

		// Check availability
		if event.TotalTickets < req.TicketCount {
			return fmt.Errorf("not enough tickets available. Only %d tickets left", event.TotalTickets)
		}

		// Decrement tickets
		if err := s.eventRepo.DecrementTickets(ctx, req.EventID, req.TicketCount); err != nil {
			return err
		}

		// Create booking
		booking = &models.Booking{
			UserID:      userID,
			EventID:     req.EventID,
			TicketCount: req.TicketCount,
			TotalPrice:  float64(req.TicketCount) * event.TicketPrice,
			Status:      models.BookingStatusPending,
			ExpiresAt:   time.Now().Add(s.timeout),
		}

		if err := s.bookingRepo.Create(ctx, booking); err != nil {
			return fmt.Errorf("failed to create booking: %w", err)
		}

		return nil
	})

	return booking, err
}

func (s *bookingService) GetBooking(ctx context.Context, id int) (*models.BookingWithDetails, error) {
	booking, err := s.bookingRepo.GetWithDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	return booking, nil
}

func (s *bookingService) GetUserBookings(ctx context.Context, userID int) ([]*models.Booking, error) {
	bookings, err := s.bookingRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user bookings: %w", err)
	}
	return bookings, nil
}

func (s *bookingService) ConfirmPayment(ctx context.Context, bookingID int) error {
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	if booking.Status != models.BookingStatusPending {
		return fmt.Errorf("booking is not in pending status")
	}

	if time.Now().After(booking.ExpiresAt) {
		return fmt.Errorf("booking has expired")
	}

	if err := s.bookingRepo.UpdateStatus(ctx, bookingID, models.BookingStatusConfirmed); err != nil {
		return fmt.Errorf("failed to confirm booking: %w", err)
	}

	return nil
}

func (s *bookingService) CancelBooking(ctx context.Context, bookingID int) error {
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	if booking.Status == models.BookingStatusCancelled {
		return fmt.Errorf("booking already cancelled")
	}

	if booking.Status == models.BookingStatusConfirmed {
		return fmt.Errorf("cannot cancel confirmed booking")
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update booking status
		if err := s.bookingRepo.UpdateStatus(ctx, bookingID, models.BookingStatusCancelled); err != nil {
			return fmt.Errorf("failed to cancel booking: %w", err)
		}

		// Release tickets back to event
		if err := s.eventRepo.IncrementTickets(ctx, booking.EventID, booking.TicketCount); err != nil {
			return fmt.Errorf("failed to release tickets: %w", err)
		}

		return nil
	})
}

func (s *bookingService) ProcessExpiredBookings(ctx context.Context) error {
	expiredBookings, err := s.bookingRepo.GetExpiredPending(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired bookings: %w", err)
	}

	for _, booking := range expiredBookings {
		if err := s.CancelBooking(ctx, booking.ID); err != nil {
			// Log error but continue processing
			fmt.Printf("Failed to cancel expired booking %d: %v\n", booking.ID, err)
		}
	}

	return nil
}
