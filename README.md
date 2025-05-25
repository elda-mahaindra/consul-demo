# Consul Demo

## Author

Elda Mahaindra ([faith030@gmail.com](mailto:faith030@gmail.com))

## Overview

This project demonstrates how to use Consul as a service registry for microservices architecture. The demo showcases service discovery patterns where services automatically register themselves with Consul, and an API gateway dynamically discovers and routes requests to these services without requiring code changes when new services are added.

## Architecture

The demo consists of the following components:

### Services

- **service-a**: A REST service that self-registers with Consul and exposes a `/ping` endpoint
- **service-b**: A REST service that self-registers with Consul and exposes a `/ping` endpoint
- **api-gateway**: A gateway service that dynamically discovers services from Consul and routes requests

### Infrastructure

- **Consul**: Service registry and discovery server

## Key Features

### Service Discovery

- Services automatically register themselves with Consul on startup
- Each service provides health check endpoints for Consul monitoring
- Services deregister gracefully on shutdown

### Dynamic Routing

- API Gateway queries Consul to discover available services
- No hardcoded service endpoints in the gateway
- Automatic load balancing when multiple instances of the same service are running

### Scalability

- New services (e.g., service-c) can be added without modifying the API Gateway
- No need to rebuild or redeploy the gateway when adding new services
- Zero-downtime service discovery

## Demo Workflow

1. **Service Registration**: service-a and service-b start up and register themselves with Consul
2. **Health Monitoring**: Consul continuously monitors service health via health check endpoints
3. **Service Discovery**: API Gateway queries Consul to get the list of available services and their endpoints
4. **Dynamic Routing**: API Gateway routes `/ping` requests to the appropriate service based on Consul registry
5. **Extensibility**: Adding service-c requires no changes to the API Gateway - it will automatically discover and route to the new service

## Benefits

- **Zero Configuration**: No manual service endpoint configuration required
- **High Availability**: Automatic failover when services become unavailable
- **Scalability**: Easy horizontal scaling of services
- **Maintainability**: Reduced coupling between services and gateway
- **Operational Efficiency**: No downtime required when adding new services

## Getting Started

[Instructions for running the demo will be added here]
