package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Simplified auth - in production use JWT
		userIDHeader := c.Get("X-User-ID")
		if userIDHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authentication",
			})
		}

		userID, err := strconv.Atoi(userIDHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}
