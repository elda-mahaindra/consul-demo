package middleware

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler creates a middleware for centralized error handling
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Forward to next handler
		err := c.Next()

		// Check if response was written
		if len(c.Response().Body()) == 0 {
			if err == nil {
				// No error but no response sent - this is a handler bug
				log.Printf("Warning: Handler didn't send any response for %s %s\n",
					c.Method(), c.Path())

				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal server error: no response sent",
				})
			}
		} else if err == nil {
			// Response was sent and no error - all good
			return nil
		}

		// Handle fiber errors
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"error": fiberErr.Message,
			})
		}

		// Default error response
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
}
