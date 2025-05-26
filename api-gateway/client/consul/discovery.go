package consul

import (
	"fmt"
	"math/rand"
	"time"

	"api-gateway/util/config"

	"github.com/hashicorp/consul/api"
)

// DiscoveryClient handles service discovery using Consul
type DiscoveryClient struct {
	client *api.Client
}

// ServiceInstance represents a discovered service instance
type ServiceInstance struct {
	ID      string
	Name    string
	Address string
	Port    int
	Tags    []string
	Meta    map[string]string
}

// NewDiscoveryClient creates a new Consul discovery client
func NewDiscoveryClient(config config.Consul) (*DiscoveryClient, error) {
	// Create Consul client configuration
	consulConfig := api.DefaultConfig()
	consulConfig.Address = fmt.Sprintf("%s:%d", config.Host, config.Port)
	consulConfig.Scheme = config.Scheme

	// Create Consul client
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &DiscoveryClient{
		client: client,
	}, nil
}

// DiscoverService finds healthy instances of a service
// Returns all available instances for load balancing
func (d *DiscoveryClient) DiscoverService(serviceName string) ([]ServiceInstance, error) {
	// Query Consul for healthy instances of the service
	services, _, err := d.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service %s: %w", serviceName, err)
	}

	// Check if any instances are available
	if len(services) == 0 {
		return nil, fmt.Errorf("no healthy instances of service %s found", serviceName)
	}

	// Convert Consul response to our ServiceInstance format
	var instances []ServiceInstance
	for _, service := range services {
		instance := ServiceInstance{
			ID:      service.Service.ID,
			Name:    service.Service.Service,
			Address: service.Service.Address,
			Port:    service.Service.Port,
			Tags:    service.Service.Tags,
			Meta:    service.Service.Meta,
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// DiscoverServiceWithLoadBalancing finds a service and returns one instance using round-robin
func (d *DiscoveryClient) DiscoverServiceWithLoadBalancing(serviceName string) (*ServiceInstance, error) {
	instances, err := d.DiscoverService(serviceName)
	if err != nil {
		return nil, err
	}

	// Simple random load balancing
	// In production, you might want more sophisticated algorithms
	rand.Seed(time.Now().UnixNano())
	selectedInstance := instances[rand.Intn(len(instances))]

	return &selectedInstance, nil
}

// GetAllServices returns all available services in Consul
func (d *DiscoveryClient) GetAllServices() (map[string][]string, error) {
	services, _, err := d.client.Catalog().Services(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get all services: %w", err)
	}

	return services, nil
}
