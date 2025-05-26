package service

import (
	"fmt"
	"log"

	"api-gateway/client/consul"
)

type PingServiceParam struct {
	ServiceName string
}

// PingServiceResponse represents the response from a ping request
type PingServiceResponse struct {
	Service     string                  `json:"service"`
	Message     string                  `json:"message"`
	Instance    *consul.ServiceInstance `json:"instance"`
	StatusCode  int                     `json:"status_code"`
	RawResponse map[string]interface{}  `json:"raw_response"`
}

// PingService discovers and pings a specific service
// This is the core function that demonstrates dynamic service discovery
func (s *Service) PingService(param *PingServiceParam) (*PingServiceResponse, error) {
	log.Printf("🔍 Discovering service: %s", param.ServiceName)

	// 1. Discover the service using Consul
	instance, err := s.discoveryClient.DiscoverServiceWithLoadBalancing(param.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("service discovery failed for %s: %w", param.ServiceName, err)
	}

	log.Printf("✅ Found service instance: %s at %s:%d", instance.Name, instance.Address, instance.Port)

	// 2. Build the URL dynamically
	url := fmt.Sprintf("http://%s:%d/ping", instance.Address, instance.Port)

	log.Printf("🌐 Making request to: %s", url)

	// 3. Make the HTTP request
	response, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to ping service %s at %s: %w", param.ServiceName, url, err)
	}

	log.Printf("📨 Received response with status: %d", response.StatusCode)

	// 4. Return structured response
	return &PingServiceResponse{
		Service:     param.ServiceName,
		Message:     fmt.Sprintf("Successfully pinged %s", param.ServiceName),
		Instance:    instance,
		StatusCode:  response.StatusCode,
		RawResponse: response.Body,
	}, nil
}
