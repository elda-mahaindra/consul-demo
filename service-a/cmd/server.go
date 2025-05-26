package main

import (
	"fmt"
	"log"
	"os"

	"service-a/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func runRestServer(port int, api *api.Api) {
	// Init fiber app
	app := fiber.New()

	// CORS middleware configuration
	corsConfig := cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}

	app.Use(cors.New(corsConfig))

	// Endpoint definitions
	app = api.DefineEndpoints(app)

	// start the server
	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("failed to listen at port: %v!", port)

		os.Exit(1)
	}

	log.Printf("rest server started successfully 🚀")
}
