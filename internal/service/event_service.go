package service

import (
	"context"
	"event-booking-be/internal/models"
	"event-booking-be/internal/repository"
	"fmt"
)

type eventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{
		eventRepo: eventRepo,
	}
}

func (s *eventService) CreateEvent(ctx context.Context, req *models.CreateEventRequest) (*models.Event, error) {
	event := &models.Event{
		Name:         req.Name,
		Description:  req.Description,
		DateTime:     req.DateTime,
		TotalTickets: req.TotalTickets,
		TicketPrice:  req.TicketPrice,
	}

	if err := s.eventRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

func (s *eventService) GetEvent(ctx context.Context, id int) (*models.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}
	return event, nil
}

func (s *eventService) GetAllEvents(ctx context.Context) ([]*models.Event, error) {
	events, err := s.eventRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	return events, nil
}

func (s *eventService) UpdateEvent(ctx context.Context, id int, req *models.UpdateEventRequest) (*models.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	if req.Name != nil {
		event.Name = *req.Name
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.DateTime != nil {
		event.DateTime = *req.DateTime
	}
	if req.TotalTickets != nil {
		event.TotalTickets = *req.TotalTickets
	}
	if req.TicketPrice != nil {
		event.TicketPrice = *req.TicketPrice
	}

	if err := s.eventRepo.Update(ctx, id, event); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return event, nil
}

func (s *eventService) DeleteEvent(ctx context.Context, id int) error {
	if err := s.eventRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}

func (s *eventService) GetEventStatistics(ctx context.Context, eventID int) (*models.EventStatistics, error) {
	stats, err := s.eventRepo.GetStatsByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event statistics: %w", err)
	}
	return stats, nil
}
