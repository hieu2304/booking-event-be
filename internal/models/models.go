package models

import (
	"time"

	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
)

type Event struct {
	ID           int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	DateTime     time.Time      `gorm:"not null" json:"date_time"`
	TotalTickets int            `gorm:"not null" json:"total_tickets"`
	TicketPrice  float64        `gorm:"type:decimal(10,2);not null" json:"ticket_price"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Bookings     []Booking      `gorm:"foreignKey:EventID" json:"-"`
}

func (Event) TableName() string {
	return "events"
}

type User struct {
	ID        int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Bookings  []Booking      `gorm:"foreignKey:UserID" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type Booking struct {
	ID          int            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int            `gorm:"not null;index" json:"user_id"`
	EventID     int            `gorm:"not null;index" json:"event_id"`
	TicketCount int            `gorm:"not null" json:"ticket_count"`
	TotalPrice  float64        `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status      BookingStatus  `gorm:"type:varchar(20);not null;index" json:"status"`
	ExpiresAt   time.Time      `gorm:"not null;index" json:"expires_at"`
	ConfirmedAt *time.Time     `json:"confirmed_at,omitempty"`
	CancelledAt *time.Time     `json:"cancelled_at,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
	Event       Event          `gorm:"foreignKey:EventID" json:"-"`
}

func (Booking) TableName() string {
	return "bookings"
}

// DTOs

type CreateEventRequest struct {
	Name         string    `json:"name" validate:"required"`
	Description  string    `json:"description"`
	DateTime     time.Time `json:"date_time" validate:"required"`
	TotalTickets int       `json:"total_tickets" validate:"required,min=1"`
	TicketPrice  float64   `json:"ticket_price" validate:"required,min=0"`
}

type UpdateEventRequest struct {
	Name         *string    `json:"name,omitempty"`
	Description  *string    `json:"description,omitempty"`
	DateTime     *time.Time `json:"date_time,omitempty"`
	TotalTickets *int       `json:"total_tickets,omitempty"`
	TicketPrice  *float64   `json:"ticket_price,omitempty"`
}

type EventStatistics struct {
	EventID        int     `json:"event_id"`
	EventName      string  `json:"event_name"`
	TotalTickets   int     `json:"total_tickets"`
	TicketsSold    int     `json:"tickets_sold"`
	TicketsLeft    int     `json:"tickets_left"`
	Revenue        float64 `json:"revenue"`
	PendingBooking int     `json:"pending_bookings"`
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type CreateBookingRequest struct {
	EventID     int `json:"event_id" validate:"required"`
	TicketCount int `json:"ticket_count" validate:"required,min=1"`
}

type BookingWithDetails struct {
	Booking
	UserName      string    `json:"user_name"`
	UserEmail     string    `json:"user_email"`
	EventName     string    `json:"event_name"`
	EventDateTime time.Time `json:"event_date_time"`
}
