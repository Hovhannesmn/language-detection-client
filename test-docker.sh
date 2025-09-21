#!/bin/bash

# Docker Test Script for Language Detection Client
set -e

echo "🐳 Testing Docker Setup for Language Detection Client"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
print_status "Checking Docker status..."
if ! docker info >/dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi
print_success "Docker is running"

# Check Docker Compose
print_status "Checking Docker Compose..."
if ! docker-compose --version >/dev/null 2>&1; then
    print_error "Docker Compose is not available."
    exit 1
fi
print_success "Docker Compose is available"

# Start the mock language detection service
print_status "Starting mock language detection service..."
docker-compose up language-detection-service -d

# Wait for service to be healthy
print_status "Waiting for service to be healthy..."
for i in {1..30}; do
    if curl -s http://localhost:6012/health >/dev/null 2>&1; then
        print_success "Mock service is healthy"
        break
    fi
    if [ $i -eq 30 ]; then
        print_error "Service failed to become healthy"
        exit 1
    fi
    sleep 1
done

# Test the mock API
print_status "Testing mock language detection API..."
response=$(curl -s -X POST http://localhost:6012/api/v1/detect)
expected='{"language_code": "en-US", "confidence": 0.95}'
if [ "$response" = "$expected" ]; then
    print_success "Mock API is working correctly"
else
    print_warning "Mock API response: $response"
fi

# Test the web interface
print_status "Testing web interface..."
if curl -s http://localhost:6012/ | grep -q "Language Detection Service"; then
    print_success "Web interface is accessible"
else
    print_warning "Web interface may not be working correctly"
fi

# Show service status
print_status "Docker container status:"
docker-compose ps

echo ""
print_success "Docker setup test completed successfully!"
echo ""
echo "🌐 Mock Language Detection Service is running at:"
echo "   - Health Check: http://localhost:6012/health"
echo "   - Web Interface: http://localhost:6012/"
echo "   - API Endpoint: http://localhost:6012/api/v1/detect"
echo ""
echo "📋 Available commands:"
echo "   - make docker-status    # Check container status"
echo "   - make docker-logs      # View service logs"
echo "   - make docker-clean     # Clean up resources"
echo "   - docker-compose down   # Stop services"
echo ""
echo "🧪 To test the client (when built):"
echo "   - make docker-compose-test     # Run automated tests"
echo "   - make docker-compose-interactive  # Start interactive client"
