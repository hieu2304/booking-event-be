package handler

import (
	"event-booking-be/internal/models"
	"event-booking-be/internal/service"
	"event-booking-be/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	eventService service.EventService
}

func NewEventHandler(eventService service.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	var req models.CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequestResponse(c, utils.INVALID_REQUEST_BODY, "Invalid request body")
	}

	event, err := h.eventService.CreateEvent(c.Context(), &req)
	if err != nil {
		return InternalErrorResponse(c, utils.EVENT_CREATE_FAILED, err.Error())
	}

	return CreatedResponse(c, event)
}

func (h *EventHandler) GetEvent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequestResponse(c, utils.EVENT_INVALID_ID, "Invalid event ID")
	}

	event, err := h.eventService.GetEvent(c.Context(), id)
	if err != nil {
		return NotFoundResponse(c, utils.EVENT_NOT_FOUND, "Event not found")
	}

	return SuccessResponse(c, event)
}

func (h *EventHandler) GetAllEvents(c *fiber.Ctx) error {
	events, err := h.eventService.GetAllEvents(c.Context())
	if err != nil {
		return InternalErrorResponse(c, utils.INTERNAL_SERVER_ERROR, err.Error())
	}

	return SuccessResponse(c, events)
}

func (h *EventHandler) UpdateEvent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequestResponse(c, utils.EVENT_INVALID_ID, "Invalid event ID")
	}

	var req models.UpdateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequestResponse(c, utils.INVALID_REQUEST_BODY, "Invalid request body")
	}

	event, err := h.eventService.UpdateEvent(c.Context(), id, &req)
	if err != nil {
		return InternalErrorResponse(c, utils.EVENT_UPDATE_FAILED, err.Error())
	}

	return SuccessResponse(c, event)
}

func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequestResponse(c, utils.EVENT_INVALID_ID, "Invalid event ID")
	}

	if err := h.eventService.DeleteEvent(c.Context(), id); err != nil {
		return InternalErrorResponse(c, utils.EVENT_DELETE_FAILED, err.Error())
	}

	return SuccessResponse(c, fiber.Map{"message": "Event deleted"})
}

func (h *EventHandler) GetEventStatistics(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequestResponse(c, utils.EVENT_INVALID_ID, "Invalid event ID")
	}

	stats, err := h.eventService.GetEventStatistics(c.Context(), id)
	if err != nil {
		return InternalErrorResponse(c, utils.EVENT_NOT_FOUND, err.Error())
	}

	return SuccessResponse(c, stats)
}
