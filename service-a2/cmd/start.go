package main

import (
	"log"
	"os"
	"os/signal"

	"service-a2/api"
	"service-a2/util/config"
)

func start() {
	// Load environment variables from .env file
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Printf("failed to load config: %v", err)

		os.Exit(1)
	}

	log.Printf("config: %v", config)
	log.Printf("Starting %s service ...", config.App.Name)

	// consul registration
	err = consulRegistration(config)
	if err != nil {
		log.Printf("failed to register service: %v", err)

		os.Exit(1)
	}

	// Init api layer
	restApi := api.NewApi(config.App.Name)

	// Run rest server
	runRestServer(config.App.Port, restApi)

	// wait for ctrl + c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block until a signal is received
	<-ch

	log.Printf("end of program...")
}
