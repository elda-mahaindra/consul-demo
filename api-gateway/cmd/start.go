package main

import (
	"log"
	"os"
	"os/signal"

	"api-gateway/api"
	"api-gateway/client/consul"
	"api-gateway/client/http_adapter"
	"api-gateway/service"
	"api-gateway/util/config"
)

func start() {
	// Load environment variables from .env file
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Printf("failed to load config: %v", err)
		os.Exit(1)
	}

	log.Printf("Starting %s service ...", config.App.Name)

	// Init Consul discovery client
	log.Printf("Initializing Consul discovery client...")
	discoveryClient, err := consul.NewDiscoveryClient(config.Consul)
	if err != nil {
		log.Printf("failed to initialize Consul discovery client: %v", err)
		os.Exit(1)
	}
	log.Printf("✅ Consul discovery client initialized successfully")

	// Init HTTP client
	httpClient := http_adapter.NewClient()

	// Init service layer with both HTTP client and discovery client
	service := service.NewService(httpClient, discoveryClient)

	// Init API layer
	restApi := api.NewApi(config.App.Name, service)

	// Run rest server
	runRestServer(config.App.Port, restApi)

	// wait for ctrl + c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block until a signal is received
	<-ch

	log.Printf("end of program...")
}
