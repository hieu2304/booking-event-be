package handler

import (
	"event-booking-be/internal/models"
	"event-booking-be/internal/service"
	"event-booking-be/internal/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	var req models.CreateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequestResponse(c, utils.INVALID_REQUEST_BODY, "Invalid request body")
	}

	booking, err := h.bookingService.CreateBooking(c.Context(), userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not enough tickets") {
			return BadRequestResponse(c, utils.BOOKING_NOT_ENOUGH_TICKETS, err.Error())
		}
		return BadRequestResponse(c, utils.BOOKING_CREATE_FAILED, err.Error())
	}

	return CreatedResponse(c, booking)
}

func (h *BookingHandler) GetBooking(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequestResponse(c, utils.BOOKING_INVALID_ID, "Invalid booking ID")
	}

	booking, err := h.bookingService.GetBooking(c.Context(), id)
	if err != nil {
		return NotFoundResponse(c, utils.BOOKING_NOT_FOUND, "Booking not found")
	}

	return SuccessResponse(c, booking)
}

func (h *BookingHandler) GetUserBookings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	bookings, err := h.bookingService.GetUserBookings(c.Context(), userID)
	if err != nil {
		return InternalErrorResponse(c, utils.INTERNAL_SERVER_ERROR, err.Error())
	}

	return SuccessResponse(c, bookings)
}

func (h *BookingHandler) ConfirmPayment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequestResponse(c, utils.BOOKING_INVALID_ID, "Invalid booking ID")
	}

	if err := h.bookingService.ConfirmPayment(c.Context(), id); err != nil {
		// Check error type
		if strings.Contains(err.Error(), "expired") {
			return BadRequestResponse(c, utils.BOOKING_EXPIRED, err.Error())
		}
		if strings.Contains(err.Error(), "not in pending") {
			return BadRequestResponse(c, utils.BOOKING_ALREADY_CONFIRMED, err.Error())
		}
		return BadRequestResponse(c, utils.BOOKING_NOT_FOUND, err.Error())
	}

	return SuccessResponse(c, fiber.Map{"message": "Payment confirmed"})
}

func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return BadRequestResponse(c, utils.BOOKING_INVALID_ID, "Invalid booking ID")
	}

	if err := h.bookingService.CancelBooking(c.Context(), id); err != nil {
		// Check error type
		if strings.Contains(err.Error(), "already cancelled") {
			return BadRequestResponse(c, utils.BOOKING_ALREADY_CANCELLED, err.Error())
		}
		if strings.Contains(err.Error(), "confirmed") {
			return BadRequestResponse(c, utils.BOOKING_ALREADY_CONFIRMED, "Cannot cancel confirmed booking")
		}
		return BadRequestResponse(c, utils.BOOKING_CANCEL_FAILED, err.Error())
	}

	return SuccessResponse(c, fiber.Map{"message": "Booking cancelled"})
}
