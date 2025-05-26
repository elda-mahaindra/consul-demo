package api

import (
	"api-gateway/service"

	"github.com/gofiber/fiber/v2"
)

// pingService is the core dynamic routing function
// It discovers the requested service and routes the ping request to it
func (api *Api) pingService(c *fiber.Ctx) error {
	serviceName := c.Params("serviceName")

	if serviceName == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "service name is required",
			"usage": "GET /api/ping/{service-name}",
			"examples": []string{
				"GET /api/ping/service-a",
				"GET /api/ping/service-b",
			},
		})
	}

	// Use service discovery to find and ping the service
	response, err := api.service.PingService(&service.PingServiceParam{ServiceName: serviceName})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "failed to ping service",
			"service": serviceName,
			"details": err.Error(),
		})
	}

	return c.JSON(response)
}
