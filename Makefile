# Lark Logger Makefile

.PHONY: help build test clean run lint fmt vet mod-tidy coverage

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build targets
build: ## Build the project
	@echo "Building lark-logger..."
	go build -o bin/lark-logger ./larklogger.go

build-cmd: ## Build the example command
	@echo "Building example command..."
	go build -o bin/example ./cmd/main.go

# Test targets
test: ## Run tests
	@echo "Running tests..."
	go test -v ./src/larklogger/...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./src/larklogger/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Code quality targets
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	go mod tidy

# Run targets
run: ## Run the example
	@echo "Running example..."
	@if [ -f .env.local ]; then \
		echo "Loading environment from .env.local..."; \
		export $$(cat .env.local | grep -v '^#' | xargs) && go run ./cmd/main.go; \
	else \
		echo "No .env.local found, running with default environment..."; \
		go run ./cmd/main.go; \
	fi

run-test: ## Run tests and example
	@echo "Running tests..."
	$(MAKE) test
	@echo "Running example..."
	$(MAKE) run

run-test-mode: ## Run example in test mode (with mock webhook)
	@echo "Running example in test mode..."
	@LARK_TEST_MODE=true LARK_WEBHOOK_URL=https://test.webhook.url go run ./cmd/main.go

# Clean targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

clean-all: clean ## Clean everything including go modules
	@echo "Cleaning go modules..."
	go clean -modcache

# Development targets
dev: ## Run in development mode (test + run)
	@echo "Development mode..."
	$(MAKE) test
	$(MAKE) run

install: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download

# CI targets
ci: fmt vet test ## Run CI pipeline
	@echo "CI pipeline completed successfully"

# Release targets
release-check: ## Check if ready for release
	@echo "Checking release readiness..."
	$(MAKE) fmt
	$(MAKE) vet
	$(MAKE) test
	@echo "Release check completed"

# Docker targets (if needed)
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t lark-logger .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm lark-logger

# Help target
.DEFAULT_GOAL := help
