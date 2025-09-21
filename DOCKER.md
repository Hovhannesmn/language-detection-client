# Docker Setup for Language Detection Client

This document explains how to use Docker and Docker Compose to run the Language Detection Client with a mock language detection service.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+

## Quick Start

### 1. Build and Run with Docker Compose

```bash
# Start all services (client + mock language detection service)
make docker-compose-up

# Or manually:
docker-compose up --build
```

### 2. Run Tests with Docker Compose

```bash
# Run automated tests
make docker-compose-test

# Or manually:
docker-compose --profile test up --build --abort-on-container-exit
```

### 3. Development Environment

```bash
# Start development environment with live reload
make docker-compose-dev

# Or manually:
docker-compose --profile dev up --build
```

## Available Services

### 1. Language Detection Service (Mock)
- **Container**: `ld-service`
- **Port**: `6011`
- **Purpose**: Mock language detection service for testing
- **Health Check**: `http://localhost:6011/health`
- **API**: `http://localhost:6011/api/v1/detect`

### 2. Language Detection Client
- **Container**: `ld-client` (test mode)
- **Purpose**: Runs validation tests against sample captions
- **Dependencies**: Language Detection Service

### 3. Interactive Client
- **Container**: `ld-client-interactive`
- **Purpose**: Interactive shell for manual testing
- **Usage**: `make docker-exec`

## Docker Profiles

### Test Profile
Runs automated tests with sample captions:
```bash
docker-compose --profile test up --build --abort-on-container-exit
```

### Interactive Profile
Starts an interactive client container:
```bash
docker-compose --profile interactive up --build -d
make docker-exec  # Access the container
```

### Development Profile
Development environment with live reload:
```bash
docker-compose --profile dev up --build
```

## Available Make Commands

### Docker Build Commands
```bash
make docker-build          # Build production image
make docker-build-dev      # Build development image
make docker-run           # Run container with help
```

### Docker Compose Commands
```bash
make docker-compose-up     # Start all services
make docker-compose-down   # Stop all services
make docker-compose-test   # Run tests
make docker-compose-dev    # Start development environment
make docker-compose-interactive  # Start interactive client
```

### Docker Management Commands
```bash
make docker-exec          # Access interactive container
make docker-logs          # View service logs
make docker-status        # Show container status
make docker-clean         # Clean up Docker resources
```

## Manual Docker Usage

### Build Image
```bash
docker build -t language-detection-client .
```

### Run Container
```bash
# Show help
docker run --rm language-detection-client -h

# Run with sample captions
docker run --rm -v $(pwd)/test-captions:/app/test-captions:ro \
  language-detection-client \
  -t_start=00:00:00 -t_end=00:00:20 -coverage=80 \
  -server=host.docker.internal:6011 \
  test-captions/sample.srt
```

### Interactive Mode
```bash
docker run -it --rm language-detection-client sh
```

## Directory Structure

```
.
├── Dockerfile              # Production image
├── Dockerfile.dev          # Development image
├── docker-compose.yml      # Multi-service setup
├── .dockerignore           # Docker ignore file
├── .air.toml              # Live reload config
├── docker/                # Docker configuration
│   ├── nginx.conf         # Mock service config
│   └── mock-responses/    # Static files
├── test-captions/         # Sample caption files
│   ├── sample.srt
│   └── multilingual.srt
└── output/               # Output directory
```

## Testing with Docker

### 1. Automated Tests
```bash
make docker-compose-test
```

### 2. Interactive Testing
```bash
# Start interactive client
make docker-compose-interactive

# Access the container
make docker-exec

# Inside the container, run tests:
./language-detection-client -t_start=00:00:00 -t_end=00:00:20 -coverage=80 -server=language-detection-service:80 test-captions/sample.srt
```

### 3. Manual Testing
```bash
# Start services
make docker-compose-up

# In another terminal, run client manually:
docker exec -it ld-client-interactive ./language-detection-client \
  -t_start=00:00:00 -t_end=00:00:20 -coverage=80 \
  -server=language-detection-service:80 \
  test-captions/sample.srt
```

## Troubleshooting

### Check Service Status
```bash
make docker-status
```

### View Logs
```bash
make docker-logs
```

### Clean Up
```bash
make docker-clean
```

### Health Checks
```bash
# Check mock service health
curl http://localhost:6011/health

# Check mock API
curl -X POST http://localhost:6011/api/v1/detect
```

## Environment Variables

- `LANGUAGE_SERVICE_URL`: URL of the language detection service
- `CGO_ENABLED`: Set to 0 for static builds

## Network Configuration

Services communicate through the `language-detection-network` bridge network:
- Client can reach service at `http://language-detection-service:80`
- External access via `http://localhost:6011`

## Volume Mounts

- `./test-captions:/app/test-captions:ro` - Read-only access to test files
- `./output:/app/output` - Output directory for results
- `.:/app` - Development mode (live code reload)

## Security

- Containers run as non-root user (`appuser:appgroup`)
- Minimal base images (Alpine Linux)
- No unnecessary packages or services
- Read-only file systems where possible
