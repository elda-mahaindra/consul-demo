package service

import (
	"fmt"
	"log"
)

// GetAllAvailableServices returns all services registered in Consul
func (s *Service) GetAllAvailableServices() (map[string][]string, error) {
	log.Printf("🔍 Discovering all available services")

	services, err := s.discoveryClient.GetAllServices()
	if err != nil {
		return nil, fmt.Errorf("failed to get all services: %w", err)
	}

	log.Printf("✅ Found %d services in registry", len(services))

	return services, nil
}
