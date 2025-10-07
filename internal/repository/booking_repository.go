package repository

import (
	"context"
	"event-booking-be/internal/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *bookingRepository) GetByID(ctx context.Context, id int) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.WithContext(ctx).First(&booking, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("booking not found")
	}
	return &booking, err
}

func (r *bookingRepository) GetByUserID(ctx context.Context, userID int) ([]*models.Booking, error) {
	var bookings []*models.Booking
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetByEventID(ctx context.Context, eventID int) ([]*models.Booking, error) {
	var bookings []*models.Booking
	err := r.db.WithContext(ctx).
		Where("event_id = ?", eventID).
		Order("created_at DESC").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id int, status models.BookingStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	
	if status == models.BookingStatusConfirmed {
		updates["confirmed_at"] = time.Now()
	} else if status == models.BookingStatusCancelled {
		updates["cancelled_at"] = time.Now()
	}
	
	result := r.db.WithContext(ctx).Model(&models.Booking{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		return fmt.Errorf("booking not found")
	}
	return result.Error
}

func (r *bookingRepository) GetExpiredPending(ctx context.Context) ([]*models.Booking, error) {
	var bookings []*models.Booking
	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at < ?", models.BookingStatusPending, time.Now()).
		Order("expires_at ASC").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetWithDetails(ctx context.Context, id int) (*models.BookingWithDetails, error) {
	var booking models.BookingWithDetails
	err := r.db.WithContext(ctx).
		Table("bookings b").
		Select(`
			b.id, b.user_id, b.event_id, b.ticket_count, b.total_price, b.status,
			b.expires_at, b.confirmed_at, b.cancelled_at, b.created_at, b.updated_at,
			u.name as user_name, u.email as user_email, e.name as event_name, e.date_time as event_date_time
		`).
		Joins("JOIN users u ON b.user_id = u.id").
		Joins("JOIN events e ON b.event_id = e.id").
		Where("b.id = ? AND b.deleted_at IS NULL", id).
		Scan(&booking).Error
	
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("booking not found")
	}
	return &booking, err
}

func (r *bookingRepository) GetStatsByEventID(ctx context.Context, eventID int) (*models.EventStatistics, error) {
	// This method moved to event_repository
	return nil, fmt.Errorf("use EventRepository.GetStatsByEventID instead")
}
