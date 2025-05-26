package api

import (
	"github.com/gofiber/fiber/v2"
)

// pingAllServices pings all available services
func (api *Api) pingAllServices(c *fiber.Ctx) error {
	results, err := api.service.PingAllServices()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "failed to ping all services",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"results": results,
		"count":   len(results),
		"message": "Ping results for all discovered services",
	})
}
