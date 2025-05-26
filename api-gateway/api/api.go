package api

import (
	"api-gateway/middleware"
	"api-gateway/service"

	"github.com/gofiber/fiber/v2"
)

type Api struct {
	serviceName string

	service *service.Service
}

func NewApi(serviceName string, service *service.Service) *Api {
	return &Api{
		serviceName: serviceName,

		service: service,
	}
}

func (api *Api) DefineEndpoints(app *fiber.App) *fiber.App {
	// Error handler middleware
	app.Use(middleware.ErrorHandler())

	// Service Discovery Routes
	discovery := app.Group("/discovery")

	// Get all available services
	discovery.Get("/services", api.getAllServices)

	// Ping all available services
	discovery.Get("/ping-all", api.pingAllServices)

	// Generic Service Routing
	// This is the main feature - dynamic routing to any service!
	routes := app.Group("/api")

	// Generic ping endpoint - routes to any service dynamically
	// Usage: GET /api/ping/{service-name}
	// Examples:
	//   GET /api/ping/service-a  -> discovers and pings service-a
	//   GET /api/ping/service-b  -> discovers and pings service-b
	//   GET /api/ping/service-c  -> discovers and pings service-c (when it exists)
	routes.Get("/ping/:serviceName", api.pingService)

	// Health check for the gateway itself
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": api.serviceName,
			"status":  "healthy",
			"message": "API Gateway is running",
		})
	})

	return app
}
