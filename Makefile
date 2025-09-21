# Language Detection Client Makefile

# Variables
BINARY_NAME=language-detection-client
MAIN_PATH=./cmd/client
BUILD_DIR=build
TEST_CAPTIONS_FILE=test_captions.srt

# Default server address
SERVER_ADDR=localhost:6011

# Build flags
BUILD_FLAGS=-ldflags="-s -w"

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build clean run run-test test lint fmt vet tidy deps install dev build-linux build-windows build-darwin

# Build targets
build: ## Build the application
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Development targets
dev: build ## Build and run in development mode
	@echo "$(BLUE)Running in development mode...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) -h

install: ## Install dependencies
	@echo "$(BLUE)Installing dependencies...$(NC)"
	@go mod tidy
	@go mod download
	@echo "$(GREEN)Dependencies installed$(NC)"

deps: install ## Alias for install

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)Vet completed$(NC)"

lint: vet ## Run linting checks

tidy: ## Tidy go modules
	@echo "$(BLUE)Tidying go modules...$(NC)"
	@go mod tidy
	@echo "$(GREEN)Modules tidied$(NC)"

# Test targets
test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)Tests completed$(NC)"

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

test-short: ## Run tests excluding integration tests
	@echo "$(BLUE)Running short tests...$(NC)"
	@go test -v -short ./...
	@echo "$(GREEN)Short tests completed$(NC)"

test-unit: ## Run only unit tests (no integration tests)
	@echo "$(BLUE)Running unit tests...$(NC)"
	@go test -v -run="^Test" ./... | grep -v "Integration"
	@echo "$(GREEN)Unit tests completed$(NC)"

test-integration: ## Run only integration tests
	@echo "$(BLUE)Running integration tests...$(NC)"
	@go test -v -run=".*Integration" ./...
	@echo "$(GREEN)Integration tests completed$(NC)"

test-parser: ## Run parser tests only
	@echo "$(BLUE)Running parser tests...$(NC)"
	@go test -v ./internal/parser
	@echo "$(GREEN)Parser tests completed$(NC)"

test-validator: ## Run validator tests only
	@echo "$(BLUE)Running validator tests...$(NC)"
	@go test -v ./internal/validator
	@echo "$(GREEN)Validator tests completed$(NC)"

test-client: ## Run client tests only
	@echo "$(BLUE)Running client tests...$(NC)"
	@go test -v ./internal/client
	@echo "$(GREEN)Client tests completed$(NC)"

test-config: ## Run config tests only
	@echo "$(BLUE)Running config tests...$(NC)"
	@go test -v ./internal/config
	@echo "$(GREEN)Config tests completed$(NC)"

test-watch: ## Run tests in watch mode (requires entr)
	@echo "$(BLUE)Running tests in watch mode...$(NC)"
	@which entr > /dev/null || (echo "$(RED)entr not found. Install with: brew install entr$(NC)" && exit 1)
	@find . -name "*.go" | entr -c make test

generate-mocks: ## Generate mock files
	@echo "$(BLUE)Generating mocks...$(NC)"
	@go generate ./...
	@echo "$(GREEN)Mocks generated$(NC)"

# Run targets with different parameters
run-help: build ## Show help message
	@$(BUILD_DIR)/$(BINARY_NAME) -h

run-test: build create-test-file ## Run with test captions file
	@echo "$(BLUE)Running with test captions...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) -t_start=00:00:00 -t_end=00:00:20 -coverage=80 -server=$(SERVER_ADDR) $(TEST_CAPTIONS_FILE)

run-coverage-test: build create-test-file ## Run coverage validation test
	@echo "$(BLUE)Running coverage validation test...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) -t_start=00:00:01 -t_end=00:00:10 -coverage=50 -server=$(SERVER_ADDR) $(TEST_CAPTIONS_FILE)

run-time-test: build create-test-file ## Run with custom time range
	@echo "$(BLUE)Running with custom time range...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) -t_start=00:00:05 -t_end=00:00:15 -coverage=70 -server=$(SERVER_ADDR) $(TEST_CAPTIONS_FILE)

# Docker-based run with parameters
run: ## Run client with configurable parameters (FILE, START, END, COVERAGE, SERVER)
	@echo "$(BLUE)Running language detection client...$(NC)"
	@echo "$(YELLOW)Parameters:$(NC)"
	@echo "  FILE: $(FILE)"
	@echo "  START: $(START)"
	@echo "  END: $(END)"
	@echo "  COVERAGE: $(COVERAGE)"
	@echo "  SERVER: $(SERVER)"
	@echo ""
	@if [ -n "$(FILE)" ] && [ ! -f "$(FILE)" ]; then \
		echo "$(RED)Error: File $(FILE) not found$(NC)"; \
		exit 1; \
	fi
	@echo "version: '3.8'" > docker-compose.override.yml
	@echo "services:" >> docker-compose.override.yml
	@echo "  language-detection-client:" >> docker-compose.override.yml
	@echo "    command:" >> docker-compose.override.yml
	@echo "      - \"./language-detection-client\"" >> docker-compose.override.yml
	@echo "      - \"-server=host.docker.internal:6011\"" >> docker-compose.override.yml
	@echo "      - \"-t_start=$(START)\"" >> docker-compose.override.yml
	@echo "      - \"-t_end=$(END)\"" >> docker-compose.override.yml
	@echo "      - \"-coverage=$(COVERAGE)\"" >> docker-compose.override.yml
	@echo "      - \"/app/$(FILE)\"" >> docker-compose.override.yml
	@SERVER=$(SERVER) START=$(START) END=$(END) COVERAGE=$(COVERAGE) docker-compose up --build
	@rm -f docker-compose.override.yml

# Utility targets
create-test-file: ## Create a test captions file
	@echo "$(BLUE)Creating test captions file...$(NC)"
	@echo "1" > $(TEST_CAPTIONS_FILE)
	@echo "00:00:01,000 --> 00:00:04,000" >> $(TEST_CAPTIONS_FILE)
	@echo "Hello, this is a test caption in English." >> $(TEST_CAPTIONS_FILE)
	@echo "" >> $(TEST_CAPTIONS_FILE)
	@echo "2" >> $(TEST_CAPTIONS_FILE)
	@echo "00:00:05,000 --> 00:00:08,000" >> $(TEST_CAPTIONS_FILE)
	@echo "Bonjour, ceci est un test en français." >> $(TEST_CAPTIONS_FILE)
	@echo "" >> $(TEST_CAPTIONS_FILE)
	@echo "3" >> $(TEST_CAPTIONS_FILE)
	@echo "00:00:09,000 --> 00:00:12,000" >> $(TEST_CAPTIONS_FILE)
	@echo "Hola, esto es una prueba en español." >> $(TEST_CAPTIONS_FILE)
	@echo "" >> $(TEST_CAPTIONS_FILE)
	@echo "4" >> $(TEST_CAPTIONS_FILE)
	@echo "00:00:13,000 --> 00:00:16,000" >> $(TEST_CAPTIONS_FILE)
	@echo "Ciao, questo è un test in italiano." >> $(TEST_CAPTIONS_FILE)
	@echo "" >> $(TEST_CAPTIONS_FILE)
	@echo "5" >> $(TEST_CAPTIONS_FILE)
	@echo "00:00:17,000 --> 00:00:20,000" >> $(TEST_CAPTIONS_FILE)
	@echo "Guten Tag, das ist ein Test auf Deutsch." >> $(TEST_CAPTIONS_FILE)
	@echo "$(GREEN)Test file created: $(TEST_CAPTIONS_FILE)$(NC)"

clean: ## Clean build artifacts and test files
	@echo "$(BLUE)Cleaning up...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f $(TEST_CAPTIONS_FILE)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Cleanup completed$(NC)"

clean-all: clean ## Clean everything including go cache
	@echo "$(BLUE)Cleaning go cache...$(NC)"
	@go clean -cache
	@echo "$(GREEN)All cleanup completed$(NC)"

# Docker targets
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	@docker build -t $(BINARY_NAME) .
	@echo "$(GREEN)Docker image built: $(BINARY_NAME)$(NC)"


docker-run: ## Run in Docker container
	@echo "$(BLUE)Running in Docker...$(NC)"
	@docker run --rm $(BINARY_NAME) -h


docker-compose-up: ## Start client with docker-compose
	@echo "$(BLUE)Starting language detection client with docker-compose...$(NC)"
	@docker-compose up --build

docker-compose-down: ## Stop client with docker-compose
	@echo "$(BLUE)Stopping services with docker-compose...$(NC)"
	@docker-compose down

docker-compose-test: ## Run tests with docker-compose
	@echo "$(BLUE)Running tests with docker-compose...$(NC)"
	@docker-compose --profile test up --build --abort-on-container-exit


docker-compose-interactive: ## Start interactive client with docker-compose
	@echo "$(BLUE)Starting interactive client...$(NC)"
	@docker-compose --profile interactive up --build -d
	@echo "$(GREEN)Interactive client started. Use 'make docker-exec' to access it.$(NC)"

docker-exec: ## Execute commands in the interactive client container
	@echo "$(BLUE)Accessing interactive client container...$(NC)"
	@docker exec -it ld-client-interactive sh

docker-logs: ## View logs from all services
	@echo "$(BLUE)Viewing service logs...$(NC)"
	@docker-compose logs -f

docker-clean: ## Clean up Docker resources
	@echo "$(BLUE)Cleaning up Docker resources...$(NC)"
	@docker-compose down -v --remove-orphans
	@docker system prune -f
	@echo "$(GREEN)Docker cleanup completed$(NC)"

docker-status: ## Show status of Docker containers
	@echo "$(BLUE)Docker container status:$(NC)"
	@docker-compose ps

# CI/CD targets
ci: generate-mocks fmt vet test build ## Run CI pipeline (generate mocks, format, vet, test, build)

# Default target
.DEFAULT_GOAL := help
