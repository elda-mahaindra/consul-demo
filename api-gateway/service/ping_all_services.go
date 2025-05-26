package service

import (
	"fmt"
	"log"
)

// PingAllServices discovers and pings all available services
func (s *Service) PingAllServices() (map[string]*PingServiceResponse, error) {
	log.Printf("🔍 Discovering and pinging all services")

	// Get all services
	services, err := s.GetAllAvailableServices()
	if err != nil {
		return nil, err
	}

	results := make(map[string]*PingServiceResponse)

	// Ping each service
	for serviceName := range services {
		// Skip consul service itself
		if serviceName == "consul" {
			continue
		}

		response, err := s.PingService(&PingServiceParam{ServiceName: serviceName})
		if err != nil {
			log.Printf("❌ Failed to ping %s: %v", serviceName, err)
			results[serviceName] = &PingServiceResponse{
				Service:    serviceName,
				Message:    fmt.Sprintf("Failed to ping %s: %v", serviceName, err),
				StatusCode: 500,
			}
		} else {
			results[serviceName] = response
		}
	}

	return results, nil
}
