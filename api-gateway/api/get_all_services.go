package api

import (
	"github.com/gofiber/fiber/v2"
)

// getAllServices returns all services available in Consul
func (api *Api) getAllServices(c *fiber.Ctx) error {
	services, err := api.service.GetAllAvailableServices()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "failed to get services",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"services": services,
		"count":    len(services),
		"message":  "Available services in Consul registry",
	})
}
