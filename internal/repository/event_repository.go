package repository

import (
	"context"
	"event-booking-be/internal/models"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event *models.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *eventRepository) GetByID(ctx context.Context, id int) (*models.Event, error) {
	var event models.Event
	err := r.db.WithContext(ctx).First(&event, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("event not found")
	}
	return &event, err
}

func (r *eventRepository) GetAll(ctx context.Context) ([]*models.Event, error) {
	var events []*models.Event
	err := r.db.WithContext(ctx).Order("date_time ASC").Find(&events).Error
	return events, err
}

func (r *eventRepository) Update(ctx context.Context, id int, event *models.Event) error {
	result := r.db.WithContext(ctx).Model(&models.Event{}).Where("id = ?", id).Updates(event)
	if result.RowsAffected == 0 {
		return fmt.Errorf("event not found")
	}
	return result.Error
}

func (r *eventRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.Event{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("event not found")
	}
	return result.Error
}

func (r *eventRepository) GetAvailableTickets(ctx context.Context, eventID int) (int, error) {
	var event models.Event
	err := r.db.WithContext(ctx).Select("total_tickets").First(&event, eventID).Error
	return event.TotalTickets, err
}

func (r *eventRepository) DecrementTickets(ctx context.Context, eventID int, count int) error {
	result := r.db.WithContext(ctx).
		Model(&models.Event{}).
		Where("id = ? AND total_tickets >= ?", eventID, count).
		UpdateColumn("total_tickets", gorm.Expr("total_tickets - ?", count))
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("not enough tickets available")
	}
	return result.Error
}

func (r *eventRepository) IncrementTickets(ctx context.Context, eventID int, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.Event{}).
		Where("id = ?", eventID).
		UpdateColumn("total_tickets", gorm.Expr("total_tickets + ?", count)).Error
}

func (r *eventRepository) GetStatsByEventID(ctx context.Context, eventID int) (*models.EventStatistics, error) {
	var stats models.EventStatistics
	
	err := r.db.WithContext(ctx).
		Table("events e").
		Select(`
			e.id as event_id,
			e.name as event_name,
			e.total_tickets,
			COALESCE(SUM(CASE WHEN b.status = ? THEN b.ticket_count ELSE 0 END), 0) as tickets_sold,
			e.total_tickets - COALESCE(SUM(CASE WHEN b.status = ? THEN b.ticket_count ELSE 0 END), 0) as tickets_left,
			COALESCE(SUM(CASE WHEN b.status = ? THEN b.total_price ELSE 0 END), 0) as revenue,
			COALESCE(COUNT(CASE WHEN b.status = ? THEN 1 END), 0) as pending_booking
		`, models.BookingStatusConfirmed, models.BookingStatusConfirmed, models.BookingStatusConfirmed, models.BookingStatusPending).
		Joins("LEFT JOIN bookings b ON e.id = b.event_id AND b.deleted_at IS NULL").
		Where("e.id = ?", eventID).
		Group("e.id, e.name, e.total_tickets").
		Scan(&stats).Error
	
	return &stats, err
}

func (r *eventRepository) LockForUpdate(ctx context.Context, eventID int) (*models.Event, error) {
	var event models.Event
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&event, eventID).Error
	
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("event not found")
	}
	return &event, err
}
