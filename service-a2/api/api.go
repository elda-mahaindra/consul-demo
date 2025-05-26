package api

import (
	"service-a2/middleware"

	"github.com/gofiber/fiber/v2"
)

type Api struct {
	serviceName string
}

func NewApi(serviceName string) *Api {
	return &Api{
		serviceName: serviceName,
	}
}

func (api *Api) DefineEndpoints(app *fiber.App) *fiber.App {
	// Error handler middleware
	app.Use(middleware.ErrorHandler())

	// Ping Routes
	ping := app.Group("/ping")
	ping.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": api.serviceName,
			"message": "pong",
		})
	})

	return app
}
