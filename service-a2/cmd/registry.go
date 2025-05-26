package main

import (
	"fmt"
	"log"

	"service-a2/util/config"

	"github.com/hashicorp/consul/api"
)

// consulRegistration registers this service instance with Consul
func consulRegistration(config config.Config) error {
	// Create Consul client configuration
	consulConfig := api.DefaultConfig()
	consulConfig.Address = fmt.Sprintf("%s:%d", config.Consul.Host, config.Consul.Port)
	consulConfig.Scheme = config.Consul.Scheme

	// Create Consul client
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return fmt.Errorf("failed to create consul client: %w", err)
	}

	// Determine the registration address
	// Use RegisterAddress if specified, otherwise fall back to Host
	registerAddr := config.App.RegisterAddress
	if registerAddr == "" {
		registerAddr = config.App.Host
	}

	// Determine the health check address
	// Use HealthCheckAddress if specified, otherwise fall back to RegisterAddress
	healthCheckAddr := config.App.HealthCheckAddress
	if healthCheckAddr == "" {
		healthCheckAddr = registerAddr
	}

	// Create service registration object
	// This is the structured data Consul stores about our service
	registration := &api.AgentServiceRegistration{
		// Service ID: Unique identifier for THIS instance
		// Format: servicename-registeraddr-port (ensures uniqueness)
		ID: fmt.Sprintf("%s-%s-%d", config.App.Name, registerAddr, config.App.Port),

		// Service Name: Logical service name (what other services will search for)
		Name: config.App.Name,

		// Network location where this service can be reached by OTHER services
		// This is NOT the bind address (0.0.0.0) but the actual reachable address
		Address: registerAddr,
		Port:    config.App.Port,

		// Tags: Metadata for service discovery filtering
		// Example: API Gateway could search for services with tag "api"
		Tags: []string{"api", "rest", "microservice"},

		// Health Check: Consul will periodically hit this endpoint
		// If it returns non-200, Consul marks this instance as unhealthy
		// Unhealthy instances are excluded from service discovery results
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/ping", healthCheckAddr, config.App.Port),
			Interval:                       "10s", // Check every 10 seconds
			Timeout:                        "3s",  // Timeout after 3 seconds
			DeregisterCriticalServiceAfter: "30s", // Remove from registry if unhealthy for 30s
		},

		// Meta: Additional key-value metadata
		Meta: map[string]string{
			"version":     "1.0.0",
			"environment": "development",
			"protocol":    "http",
		},
	}

	// Register the service with Consul
	// After this call, other services can discover this service by querying Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service with consul: %w", err)
	}

	log.Printf("✅ Service '%s' successfully registered with Consul", config.App.Name)
	log.Printf("   - Service ID: %s", registration.ID)
	log.Printf("   - Bind Address: %s:%d (where service listens)", config.App.Host, config.App.Port)
	log.Printf("   - Register Address: %s:%d (where others can reach it)", registration.Address, registration.Port)
	log.Printf("   - Health Check Address: %s:%d (where Consul checks health)", healthCheckAddr, config.App.Port)
	log.Printf("   - Health Check URL: %s", registration.Check.HTTP)
	log.Printf("   - Tags: %v", registration.Tags)

	return nil
}
