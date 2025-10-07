package handler

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Success: true,
		Data:    data,
	})
}

func CreatedResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, status int, code string, message string) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Code:    code,
		Error:   message,
	})
}

func BadRequestResponse(c *fiber.Ctx, code string, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, code, message)
}

func NotFoundResponse(c *fiber.Ctx, code string, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, code, message)
}

func InternalErrorResponse(c *fiber.Ctx, code string, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, code, message)
}
