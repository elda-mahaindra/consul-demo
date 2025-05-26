# Consul Demo - Makefile
# Author: Elda Mahaindra (faith030@gmail.com)
#
# This Makefile provides convenient commands for managing the Consul service discovery demo.
# It includes commands for Docker operations, testing, monitoring, and troubleshooting.

.PHONY: help build up down restart logs status clean test-discovery test-load-balancing monitor consul-ui logs-service-a2 debug-service-a2

# Default target
help: ## Show this help message
	@echo "Consul Demo - Available Commands:"
	@echo "=================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

# Docker Management Commands
build: ## Build all Docker images
	@echo "🔨 Building all Docker images..."
	docker compose build

up: ## Start all services
	@echo "🚀 Starting all services..."
	docker compose up -d
	@echo "✅ All services started!"
	@echo "📊 Consul UI: http://localhost:8500"
	@echo "🌐 API Gateway: http://localhost:4000"

down: ## Stop all services
	@echo "🛑 Stopping all services..."
	docker compose down -v
	@echo "✅ All services stopped and volumes removed!"

restart: ## Restart all services
	@echo "🔄 Restarting all services..."
	docker compose down -v
	docker compose up -d
	@echo "✅ All services restarted with fresh volumes!"

rebuild: ## Rebuild and restart all services
	@echo "🔨 Rebuilding and restarting all services..."
	docker compose down -v
	docker compose up -d --build
	@echo "✅ All services rebuilt and restarted with fresh volumes!"

# Service-specific commands
start-consul: ## Start only Consul
	@echo "🏛️ Starting Consul..."
	docker compose up -d consul

start-services: ## Start only the microservices (service-a, service-a2, service-b)
	@echo "⚙️ Starting microservices..."
	docker compose up -d service-a service-a2 service-b

start-gateway: ## Start only the API Gateway
	@echo "🌐 Starting API Gateway..."
	docker compose up -d api-gateway

# Monitoring Commands
status: ## Show status of all containers
	@echo "📊 Container Status:"
	@echo "==================="
	docker compose ps

logs: ## Show logs for all services
	@echo "📋 Showing logs for all services..."
	docker compose logs -f

logs-consul: ## Show Consul logs
	@echo "📋 Consul logs:"
	docker compose logs -f consul

logs-service-a: ## Show service-a logs
	@echo "📋 service-a logs:"
	docker compose logs -f service-a

logs-service-a2: ## Show service-a2 logs
	@echo "📋 service-a2 logs:"
	docker compose logs -f service-a2

logs-service-b: ## Show service-b logs
	@echo "📋 service-b logs:"
	docker compose logs -f service-b

logs-gateway: ## Show API Gateway logs
	@echo "📋 API Gateway logs:"
	docker compose logs -f api-gateway

# Testing Commands
test-services: ## Test all services directly
	@echo "🧪 Testing all services directly..."
	@echo "Testing service-a (port 4001):"
	@curl -s http://localhost:4001/ping | jq '.' || echo "❌ service-a not responding"
	@echo "\nTesting service-a2 (port 4003):"
	@curl -s http://localhost:4003/ping | jq '.' || echo "❌ service-a2 not responding"
	@echo "\nTesting service-b (port 4002):"
	@curl -s http://localhost:4002/ping | jq '.' || echo "❌ service-b not responding"

test-discovery: ## Test service discovery via API Gateway
	@echo "🔍 Testing service discovery..."
	@echo "Available services:"
	@curl -s http://localhost:4000/discovery/services | jq '.' || echo "❌ API Gateway not responding"

test-routing: ## Test dynamic routing via API Gateway
	@echo "🌐 Testing dynamic routing..."
	@echo "Pinging service-a via gateway:"
	@curl -s http://localhost:4000/api/ping/service-a | jq '.' || echo "❌ Routing to service-a failed"
	@echo "\nPinging service-b via gateway:"
	@curl -s http://localhost:4000/api/ping/service-b | jq '.' || echo "❌ Routing to service-b failed"

test-load-balancing: ## Test load balancing with multiple requests
	@echo "⚖️ Testing load balancing (10 requests to service-a)..."
	@echo "Instance distribution:"
	@for i in $$(seq 1 10); do \
		curl -s http://localhost:4000/api/ping/service-a | jq -r '.instance.Address' 2>/dev/null || echo "error"; \
	done | sort | uniq -c | awk '{print "  " $$2 ": " $$1 " requests"}'

test-all: test-services test-discovery test-routing test-load-balancing ## Run all tests

# Consul-specific commands
consul-ui: ## Open Consul UI in browser
	@echo "🌐 Opening Consul UI..."
	@which open >/dev/null 2>&1 && open http://localhost:8500 || echo "Open http://localhost:8500 in your browser"

consul-members: ## Show Consul cluster members
	@echo "👥 Consul cluster members:"
	@curl -s http://localhost:8500/v1/status/leader | jq '.' || echo "❌ Consul not responding"

consul-services: ## List all registered services in Consul
	@echo "📋 Registered services in Consul:"
	@curl -s http://localhost:8500/v1/catalog/services | jq '.' || echo "❌ Consul not responding"

consul-service-a: ## Show all service-a instances
	@echo "🔍 service-a instances in Consul:"
	@curl -s http://localhost:8500/v1/health/service/service-a | jq '.[].Service | {ID, Address, Port, Tags}' || echo "❌ No service-a instances found"

consul-health: ## Check health of all services
	@echo "🏥 Health status of all services:"
	@curl -s http://localhost:8500/v1/health/state/any | jq '.[] | select(.ServiceName != "") | {Service: .ServiceName, Status: .Status, Output: .Output}' || echo "❌ Consul not responding"

# Development Commands
dev-setup: ## Setup development environment
	@echo "🛠️ Setting up development environment..."
	@echo "Checking dependencies..."
	@which docker >/dev/null 2>&1 || (echo "❌ Docker not found. Please install Docker." && exit 1)
	@which docker-compose >/dev/null 2>&1 || which docker >/dev/null 2>&1 || (echo "❌ Docker Compose not found. Please install Docker Compose." && exit 1)
	@which curl >/dev/null 2>&1 || (echo "❌ curl not found. Please install curl." && exit 1)
	@which jq >/dev/null 2>&1 || (echo "⚠️ jq not found. Install jq for better JSON output: sudo apt-get install jq" && exit 0)
	@echo "✅ All dependencies found!"

# Cleanup Commands
clean: ## Remove all containers, images, and volumes
	@echo "🧹 Cleaning up..."
	docker compose down -v --remove-orphans
	docker system prune -f
	@echo "✅ Cleanup complete!"

clean-volumes: ## Stop services and remove volumes (same as 'down')
	@echo "🗑️ Stopping services and removing volumes..."
	docker compose down -v
	@echo "✅ Services stopped and volumes removed!"

# Troubleshooting Commands
debug-service-a: ## Debug service-a issues
	@echo "🔧 Debugging service-a..."
	@echo "Container status:"
	@docker compose ps service-a
	@echo "\nRecent logs:"
	@docker compose logs --tail=20 service-a
	@echo "\nHealth check:"
	@curl -s http://localhost:4001/ping || echo "❌ Direct ping failed"

debug-service-a2: ## Debug service-a2 issues
	@echo "🔧 Debugging service-a2..."
	@echo "Container status:"
	@docker compose ps service-a2
	@echo "\nRecent logs:"
	@docker compose logs --tail=20 service-a2
	@echo "\nHealth check:"
	@curl -s http://localhost:4003/ping || echo "❌ Direct ping failed"

debug-consul: ## Debug Consul issues
	@echo "🔧 Debugging Consul..."
	@echo "Container status:"
	@docker compose ps consul
	@echo "\nRecent logs:"
	@docker compose logs --tail=20 consul
	@echo "\nConsul leader:"
	@curl -s http://localhost:8500/v1/status/leader || echo "❌ Consul API not responding"

debug-gateway: ## Debug API Gateway issues
	@echo "🔧 Debugging API Gateway..."
	@echo "Container status:"
	@docker compose ps api-gateway
	@echo "\nRecent logs:"
	@docker compose logs --tail=20 api-gateway
	@echo "\nGateway health:"
	@curl -s http://localhost:4000/health || echo "❌ Gateway health check failed"

# Quick start commands
quick-start: dev-setup build up test-all ## Complete setup: check deps, build, start, and test
	@echo "🎉 Quick start complete! All services are running and tested."

demo: up ## Start demo and show key endpoints
	@echo "🎬 Consul Demo is running!"
	@echo "=========================="
	@echo "🏛️ Consul UI:        http://localhost:8500"
	@echo "🌐 API Gateway:      http://localhost:4000"
	@echo "⚙️ service-a:        http://localhost:4001"
	@echo "⚙️ service-a2:       http://localhost:4003"
	@echo "⚙️ service-b:        http://localhost:4002"
	@echo ""
	@echo "🧪 Try these commands:"
	@echo "  make test-load-balancing  # See load balancing in action"
	@echo "  make consul-ui           # Open Consul web interface"
	@echo "  make test-all            # Run all tests"

# Port checking
check-ports: ## Check if required ports are available
	@echo "🔍 Checking required ports..."
	@for port in 4000 4001 4002 4003 8500; do \
		if lsof -i :$$port >/dev/null 2>&1; then \
			echo "❌ Port $$port is already in use"; \
		else \
			echo "✅ Port $$port is available"; \
		fi; \
	done 