# Consul Demo

## Author

Elda Mahaindra ([faith030@gmail.com](mailto:faith030@gmail.com))

## Overview

This project demonstrates how to use Consul as a service registry for microservices architecture. The demo showcases service discovery patterns where services automatically register themselves with Consul, and an API gateway dynamically discovers and routes requests to these services without requiring code changes when new services are added.

## Project Structure

```
consul-demo/
├── service-a/          # First instance of service-a (port 4001)
├── service-a2/         # Second instance of service-a (port 4003) - Load balancing demo
├── service-b/          # Independent service (port 4002)
├── api-gateway/        # Gateway with service discovery (port 4000)
├── docker-compose.yml  # Container orchestration
├── Makefile           # Convenient commands for managing the demo
└── README.md          # This file
```

**Note**: `service-a2` is a copy of `service-a` with modified configuration to demonstrate horizontal scaling and load balancing. Both instances register with the same service name but different IDs.

## Architecture

The demo consists of the following components:

### Services

- **service-a**: A REST service that self-registers with Consul and exposes a `/ping` endpoint (port 4001)
- **service-a2**: A second instance of service-a for load balancing demonstration (port 4003)
- **service-b**: A REST service that self-registers with Consul and exposes a `/ping` endpoint (port 4002)
- **api-gateway**: A gateway service that dynamically discovers services from Consul and routes requests (port 4000)

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

## Understanding Consul Service Registry

### What is Consul Service Registry?

Consul's service registry is **NOT** a simple key-value store. It's a **service catalog** that maintains a structured database of all services in your infrastructure.

### Key Concepts

#### Service vs Service Instance

- **Service**: A logical name (e.g., "user-service", "payment-service")
- **Service Instance**: A specific running copy of that service (e.g., user-service running on server1:8080)

#### Registration Data Structure

When a service registers with Consul, it provides:

```json
{
  "ID": "service-a-service-a-4001",
  "Name": "service-a",
  "Address": "service-a",
  "Port": 4001,
  "Tags": ["api", "rest", "microservice"],
  "Check": {
    "HTTP": "http://service-a:4001/ping",
    "Interval": "10s",
    "Timeout": "3s"
  },
  "Meta": {
    "version": "1.0.0",
    "environment": "development"
  }
}
```

#### How Service Discovery Works

1. **Registration Process**:

   - Service starts up
   - Service calls Consul API: "Register me with these details"
   - Consul stores the service information
   - Consul starts health checking the service

2. **Discovery Process**:
   - API Gateway needs to call "service-a"
   - Gateway queries Consul: "Give me all healthy instances of 'service-a'"
   - Consul returns healthy instances with their addresses and ports
   - Gateway picks one (load balancing) and makes the request

#### Health Checking

Consul continuously monitors service health:

- **Healthy**: Service responds to health check → included in discovery results
- **Unhealthy**: Service fails health check → excluded from discovery results
- **Critical**: Service fails for too long → automatically deregistered

## Getting Started

### Prerequisites

- Docker and Docker Compose installed
- Ports 4000, 4001, 4002, 4003, and 8500 available
- `make` utility (usually pre-installed on Linux/macOS)
- `curl` and `jq` for testing commands (optional but recommended)

### Quick Start (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd consul-demo

# One command to rule them all!
make quick-start
```

This will:

1. ✅ Check all dependencies
2. 🔨 Build all Docker images
3. 🚀 Start all services
4. 🧪 Run comprehensive tests
5. 📊 Show you everything is working

### Alternative: Step by Step

### Running the Demo

#### Quick Start with Makefile

```bash
# Check dependencies and start everything
make quick-start

# Or just start the demo
make demo

# See all available commands
make help
```

#### Manual Docker Commands

If you prefer using Docker Compose directly:

```bash
# From the root directory (consul-demo/)
docker compose up -d
```

This will start:

- **Consul server** on `http://localhost:8500`
- **service-a** on `http://localhost:4001`
- **service-a2** on `http://localhost:4003` (second instance for load balancing demo)
- **service-b** on `http://localhost:4002`
- **api-gateway** on `http://localhost:4000`

#### Verify Services are Running

```bash
# Using Makefile (recommended)
make status              # Check container status
make test-all           # Run all tests
make consul-ui          # Open Consul web interface

# Using Docker Compose directly
docker compose ps
curl http://localhost:4001/ping
curl http://localhost:4003/ping  # service-a2
curl http://localhost:4002/ping
```

#### Test Service Discovery and Load Balancing

```bash
# Using Makefile (recommended)
make test-discovery      # Check registered services
make test-routing        # Test dynamic routing
make test-load-balancing # See load balancing in action

# Using curl directly
curl http://localhost:4000/discovery/services
curl http://localhost:4000/api/ping/service-a
curl http://localhost:4000/api/ping/service-b
```

### Configuration for Docker

The services use container names for inter-container communication. Here are the configurations for both service-a instances:

**service-a (First Instance)**:

```json
{
  "app": {
    "name": "service-a",
    "host": "0.0.0.0",
    "port": 4001,
    "register_address": "service-a",
    "health_check_address": "service-a"
  },
  "consul": {
    "host": "consul",
    "port": 8500,
    "scheme": "http"
  }
}
```

**service-a2 (Second Instance)**:

```json
{
  "app": {
    "name": "service-a",
    "host": "0.0.0.0",
    "port": 4003,
    "register_address": "service-a2",
    "health_check_address": "service-a2"
  },
  "consul": {
    "host": "consul",
    "port": 8500,
    "scheme": "http"
  }
}
```

**Key Configuration Points**:

- `name`: **Same for both instances** (`"service-a"`) - this enables load balancing
- `host: "0.0.0.0"`: Bind address (where service listens on all interfaces)
- `port`: **Different ports** (4001 vs 4003) to avoid conflicts
- `register_address`: **Different container names** (service-a vs service-a2) for unique registration
- `health_check_address`: **Different addresses** for health check endpoints
- `consul.host: "consul"`: Consul container name

## API Gateway Endpoints

### Dynamic Service Routing

**Endpoint**: `GET /api/ping/{service-name}`

Route to any service by name without hardcoding URLs:

```bash
# Ping service-a
curl http://localhost:4000/api/ping/service-a

# Ping service-b
curl http://localhost:4000/api/ping/service-b

# Ping any future service (e.g., service-c)
curl http://localhost:4000/api/ping/service-c
```

**Response Example** (Load Balancing in Action):

When you call `/api/ping/service-a`, you might get either instance:

**Instance 1 Response**:

```json
{
  "service": "service-a",
  "message": "Successfully pinged service-a",
  "instance": {
    "ID": "service-a-service-a-4001",
    "Name": "service-a",
    "Address": "service-a",
    "Port": 4001,
    "Tags": ["api", "rest", "microservice"],
    "Meta": {
      "version": "1.0.0",
      "environment": "development"
    }
  },
  "status_code": 200,
  "raw_response": {
    "service": "service-a",
    "message": "pong"
  }
}
```

**Instance 2 Response**:

```json
{
  "service": "service-a",
  "message": "Successfully pinged service-a",
  "instance": {
    "ID": "service-a-service-a2-4003",
    "Name": "service-a",
    "Address": "service-a2",
    "Port": 4003,
    "Tags": ["api", "rest", "microservice"],
    "Meta": {
      "version": "1.0.0",
      "environment": "development"
    }
  },
  "status_code": 200,
  "raw_response": {
    "service": "service-a",
    "message": "pong"
  }
}
```

Notice how the `ID`, `Address`, and `Port` differ, but the `Name` is the same!

### Service Discovery

**Get All Services**: `GET /discovery/services`

```bash
curl http://localhost:4000/discovery/services
```

**Ping All Services**: `GET /discovery/ping-all`

```bash
curl http://localhost:4000/discovery/ping-all
```

## Demo Workflow

1. **Service Registration**: service-a, service-a2, and service-b start up and register themselves with Consul
   - service-a and service-a2 both register with the same service name (`"service-a"`)
   - Each gets a unique ID based on their container name and port
2. **Health Monitoring**: Consul continuously monitors service health via health check endpoints
   - Checks `http://service-a:4001/ping` for first instance
   - Checks `http://service-a2:4003/ping` for second instance
3. **Service Discovery**: API Gateway queries Consul to get the list of available services and their endpoints
   - When requesting "service-a", Consul returns both healthy instances
4. **Dynamic Routing & Load Balancing**: API Gateway routes `/ping` requests using random load balancing
   - Picks one of the available service-a instances for each request
   - No configuration needed - automatic distribution
5. **Extensibility**: Adding service-c or more service-a instances requires no changes to the API Gateway

## Benefits

- **Zero Configuration**: No manual service endpoint configuration required
- **High Availability**: Automatic failover when services become unavailable
- **Scalability**: Easy horizontal scaling of services
- **Maintainability**: Reduced coupling between services and gateway
- **Operational Efficiency**: No downtime required when adding new services

## Makefile Commands

The project includes a comprehensive Makefile with convenient commands for all operations:

### Essential Commands

```bash
make help                # Show all available commands
make demo               # Start demo and show key endpoints
make quick-start        # Complete setup: check deps, build, start, and test
make status             # Show status of all containers
make test-all           # Run all tests
```

### Docker Management

```bash
make build              # Build all Docker images
make up                 # Start all services
make down               # Stop all services (removes volumes for clean state)
make restart            # Restart all services (with fresh volumes)
make rebuild            # Rebuild and restart all services (with fresh volumes)
```

**Note**: All stop/restart commands include `-v` to remove volumes, ensuring each demo run starts with a clean state. This is ideal for demo/educational purposes.

### Testing Commands

```bash
make test-services      # Test all services directly
make test-discovery     # Test service discovery via API Gateway
make test-routing       # Test dynamic routing via API Gateway
make test-load-balancing # Test load balancing with multiple requests
```

### Monitoring and Debugging

```bash
make logs               # Show logs for all services
make logs-service-a     # Show service-a logs
make logs-service-a2    # Show service-a2 logs
make consul-ui          # Open Consul UI in browser
make consul-service-a   # Show all service-a instances
make debug-service-a2   # Debug service-a2 issues
```

### Consul-Specific Commands

```bash
make consul-services    # List all registered services
make consul-health      # Check health of all services
make consul-members     # Show Consul cluster members
```

### Utility Commands

```bash
make check-ports        # Check if required ports are available
make dev-setup          # Setup development environment
make clean              # Remove all containers, images, and volumes
```

## Docker Commands (Alternative)

If you prefer using Docker Compose directly:

### Managing Services

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Stop all services (removes volumes for clean state)
docker compose down -v

# Rebuild and restart
docker compose up -d --build
```

### Monitoring

```bash
# Check container status
docker compose ps

# View specific service logs
docker compose logs -f service-a
docker compose logs -f service-a2

# Check Consul health
curl http://localhost:8500/v1/status/leader

# View all service-a instances in Consul
curl -s http://localhost:8500/v1/health/service/service-a | jq '.[].Service | {ID, Address, Port}'
```

## Load Balancing with Multiple Instances

### Horizontal Scaling Demo

The demo includes **two instances of service-a** to demonstrate automatic load balancing:

- **service-a** (port 4001) - First instance
- **service-a2** (port 4003) - Second instance

Both register with the **same service name** (`"service-a"`) but **different IDs**:

- `service-a-service-a-4001`
- `service-a-service-a2-4003`

### Testing Load Balancing

```bash
# Make multiple requests to see load balancing in action
for i in {1..10}; do
  echo "Request $i:"
  curl -s http://localhost:4000/api/ping/service-a | jq '.instance | {ID, Address, Port}'
  echo
done
```

**Sample Output**:

```
Request 1: service-a2:4003
Request 2: service-a:4001
Request 3: service-a:4001
Request 4: service-a2:4003
...
```

### How It Works

1. **Registration**: Both instances register with Consul using the same service name
2. **Discovery**: API Gateway queries Consul for "service-a" instances
3. **Load Balancing**: Consul returns both healthy instances, Gateway picks one randomly
4. **Zero Configuration**: No code changes needed in the API Gateway!

### Technical Implementation Details

**service-a vs service-a2 Differences**:

| Aspect                  | service-a                    | service-a2                    |
| ----------------------- | ---------------------------- | ----------------------------- |
| **Go Module**           | `module service-a`           | `module service-a2`           |
| **Container Name**      | `service-a`                  | `service-a2`                  |
| **Port**                | 4001                         | 4003                          |
| **Register Address**    | `service-a`                  | `service-a2`                  |
| **Health Check URL**    | `http://service-a:4001/ping` | `http://service-a2:4003/ping` |
| **Consul Service ID**   | `service-a-service-a-4001`   | `service-a-service-a2-4003`   |
| **Consul Service Name** | `service-a` ✅ **SAME**      | `service-a` ✅ **SAME**       |

**Key Insight**: The **Service Name** is identical, enabling load balancing, while all other identifiers are unique to prevent conflicts.

## Adding New Services

To add a new service (e.g., service-c):

1. Create the service following the same pattern as service-a/service-b
2. Ensure it registers with Consul on startup
3. Add it to `docker-compose.yml`
4. Start it with `docker compose up -d service-c`
5. The API Gateway will automatically discover it - no code changes needed!

The new service will immediately be available via:

```bash
curl http://localhost:4000/api/ping/service-c
```

This demonstrates the power of service discovery - true zero-configuration service addition!

## Troubleshooting

### Common Issues with service-a2

**Issue**: service-a2 container keeps restarting

```bash
# Using Makefile (recommended)
make debug-service-a2    # Complete debugging info
make logs-service-a2     # View logs
make check-ports         # Check port conflicts

# Using Docker Compose directly
docker compose logs service-a2

# Common causes:
# 1. Port conflict (ensure 4003 is available)
# 2. Config file not found (check volume mounting)
# 3. Module import errors (ensure go.mod is correct)
```

**Issue**: service-a2 not appearing in load balancing

```bash
# Using Makefile (recommended)
make consul-service-a    # Check registered instances
make test-services       # Test direct service access

# Using curl directly
curl -s http://localhost:8500/v1/health/service/service-a | jq length
# Should return 2 (both instances)
# If only 1, check service-a2 health:
curl http://localhost:4003/ping
```

**Issue**: Load balancing not working (always same instance)

```bash
# Using Makefile (recommended)
make test-load-balancing # See distribution automatically

# Using curl directly
for i in {1..20}; do curl -s http://localhost:4000/api/ping/service-a | jq -r '.instance.Address'; done | sort | uniq -c
```

### Verifying Load Balancing

```bash
# Using Makefile (recommended)
make test-load-balancing

# Manual verification
echo "Testing load balancing..."
for i in {1..10}; do
  INSTANCE=$(curl -s http://localhost:4000/api/ping/service-a | jq -r '.instance.Address')
  echo "Request $i: $INSTANCE"
done
```

## License

This project is for educational purposes. Feel free to use it as a reference for implementing Consul service discovery in your own microservices applications.

---

_Built with ❤️ to help developers understand service discovery and dynamic routing with Consul_
