package handler

import (
	"event-booking-be/internal/models"
	"event-booking-be/internal/service"
	"event-booking-be/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequestResponse(c, utils.INVALID_REQUEST_BODY, "Invalid request body")
	}

	user, err := h.userService.CreateUser(c.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return BadRequestResponse(c, utils.USER_ALREADY_EXISTS, "Email already exists")
		}
		return InternalErrorResponse(c, utils.INTERNAL_SERVER_ERROR, err.Error())
	}

	return CreatedResponse(c, user)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return BadRequestResponse(c, utils.INVALID_REQUEST_BODY, "Invalid request body")
	}

	user, err := h.userService.Login(c.Context(), &req)
	if err != nil {
		return BadRequestResponse(c, utils.USER_INVALID_CREDENTIALS, "Invalid credentials")
	}

	return SuccessResponse(c, fiber.Map{
		"user": user,
		"message": "Login successful",
	})
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	user, err := h.userService.GetUser(c.Context(), userID)
	if err != nil {
		return NotFoundResponse(c, utils.USER_NOT_FOUND, "User not found")
	}

	return SuccessResponse(c, user)
}
