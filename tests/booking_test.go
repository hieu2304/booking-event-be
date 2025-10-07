package tests

import (
	"context"
	"event-booking-be/internal/models"
	"event-booking-be/internal/repository"
	"event-booking-be/internal/service"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Event{}, &models.User{}, &models.Booking{})
	require.NoError(t, err)

	return db
}

// Test 1: Create booking successfully
func TestCreateBooking_Success(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	eventRepo := repository.NewEventRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	userRepo := repository.NewUserRepository(db)
	bookingService := service.NewBookingService(bookingRepo, eventRepo, db, 15)

	// Setup
	event := &models.Event{
		Name:         "Concert",
		DateTime:     time.Now().Add(24 * time.Hour),
		TotalTickets: 100,
		TicketPrice:  50.0,
	}
	require.NoError(t, eventRepo.Create(ctx, event))

	user := &models.User{Name: "John", Email: "john@test.com"}
	require.NoError(t, userRepo.Create(ctx, user))

	// Test
	req := &models.CreateBookingRequest{
		EventID:     event.ID,
		TicketCount: 2,
	}
	booking, err := bookingService.CreateBooking(ctx, user.ID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, booking)
	assert.Equal(t, 2, booking.TicketCount)
	assert.Equal(t, 100.0, booking.TotalPrice)
}

// Test 2: Cannot book more tickets than available  
func TestCreateBooking_InsufficientTickets(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	eventRepo := repository.NewEventRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	userRepo := repository.NewUserRepository(db)
	bookingService := service.NewBookingService(bookingRepo, eventRepo, db, 15)

	event := &models.Event{
		Name:         "Small Event",
		DateTime:     time.Now().Add(24 * time.Hour),
		TotalTickets: 5,
		TicketPrice:  30.0,
	}
	eventRepo.Create(ctx, event)

	user := &models.User{Name: "Jane", Email: "jane@test.com"}
	userRepo.Create(ctx, user)

	req := &models.CreateBookingRequest{
		EventID:     event.ID,
		TicketCount: 10,
	}

	booking, err := bookingService.CreateBooking(ctx, user.ID, req)

	assert.Error(t, err)
	assert.Nil(t, booking)
}

// Test 3: Cancel pending booking
func TestCancelBooking_Success(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	eventRepo := repository.NewEventRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	userRepo := repository.NewUserRepository(db)
	bookingService := service.NewBookingService(bookingRepo, eventRepo, db, 15)

	event := &models.Event{
		Name:         "Event",
		DateTime:     time.Now().Add(48 * time.Hour),
		TotalTickets: 50,
		TicketPrice:  25.0,
	}
	eventRepo.Create(ctx, event)

	user := &models.User{Name: "Bob", Email: "bob@test.com"}
	userRepo.Create(ctx, user)

	req := &models.CreateBookingRequest{
		EventID:     event.ID,
		TicketCount: 3,
	}
	booking, _ := bookingService.CreateBooking(ctx, user.ID, req)

	err := bookingService.CancelBooking(ctx, booking.ID)

	assert.NoError(t, err)
	
	cancelled, _ := bookingRepo.GetByID(ctx, booking.ID)
	assert.Equal(t, models.BookingStatusCancelled, cancelled.Status)
}

// Test 4: Concurrent booking test
func TestConcurrentBookings(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	eventRepo := repository.NewEventRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	userRepo := repository.NewUserRepository(db)
	bookingService := service.NewBookingService(bookingRepo, eventRepo, db, 15)

	event := &models.Event{
		Name:         "Limited Event",
		DateTime:     time.Now().Add(48 * time.Hour),
		TotalTickets: 10,
		TicketPrice:  100.0,
	}
	eventRepo.Create(ctx, event)

	// Create users
	numUsers := 15
	users := make([]*models.User, numUsers)
	for i := 0; i < numUsers; i++ {
		user := &models.User{
			Name:  fmt.Sprintf("User%d", i),
			Email: fmt.Sprintf("user%d@test.com", i),
		}
		userRepo.Create(ctx, user)
		users[i] = user
	}

	// Concurrent bookings
	var wg sync.WaitGroup
	results := make([]error, numUsers)

	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			req := &models.CreateBookingRequest{
				EventID:     event.ID,
				TicketCount: 1,
			}
			_, err := bookingService.CreateBooking(ctx, users[idx].ID, req)
			results[idx] = err
		}(i)
	}

	wg.Wait()

	// Count successes
	successCount := 0
	for _, err := range results {
		if err == nil {
			successCount++
		}
	}

	// At most 10 should succeed
	assert.LessOrEqual(t, successCount, 10, "Cannot book more than available tickets")
	assert.Greater(t, successCount, 0, "At least some bookings should succeed")
}

// Test 5: Create event (bonus test)
func TestCreateEvent(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo)

	req := &models.CreateEventRequest{
		Name:         "Music Festival",
		Description:  "Summer music",
		DateTime:     time.Now().Add(30 * 24 * time.Hour),
		TotalTickets: 5000,
		TicketPrice:  150.0,
	}

	event, err := eventService.CreateEvent(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, "Music Festival", event.Name)
	assert.Equal(t, 5000, event.TotalTickets)
}

