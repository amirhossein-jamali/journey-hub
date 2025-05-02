.PHONY: all build clean run test lint swagger dev help mock migration-up migration-down encrypt-config test-config-reload setup-encrypted-configs

# Default target
all: lint test build

# Build the application
build:
	@echo "Building Journey Hub..."
	@go build -o bin/journeyhub ./cmd/api

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/journeyhub

# Run the application
run:
	@echo "Running Journey Hub..."
	@go run ./cmd/api

# Development mode with hot reload using air
dev:
	@echo "Starting development server with hot reload..."
	@air -c .air.toml

# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./... -cover

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage report..."
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Lint the code
lint:
	@echo "Linting code..."
	@golangci-lint run ./...

# Check for security vulnerabilities
vuln-check:
	@echo "Checking for security vulnerabilities..."
	@./scripts/govulncheck.sh

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@./scripts/init-swagger.sh

# Generate mocks for testing
mock:
	@echo "Generating mocks..."
	@go generate ./...

# Run database migrations up
migration-up:
	@echo "Running database migrations up..."
	@go run ./cmd/migration up

# Run database migrations down
migration-down:
	@echo "Running database migrations down..."
	@go run ./cmd/migration down

# Install git hooks
install-hooks:
	@echo "Installing git hooks..."
	@./scripts/install-hooks.sh

# Configuration tools
encrypt-config:
	@echo "Encrypting configuration value..."
	@go run cmd/encrypt/main.go -i

test-config-reload:
	@echo "Testing configuration hot reload..."
	@./scripts/test-config-hot-reload.sh

setup-encrypted-configs:
	@echo "Setting up encrypted configuration values for all environments..."
	@./scripts/setup-encrypted-configs.sh

# Help information
help:
	@echo "Available commands:"
	@echo "  make build           - Build the application"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make run             - Run the application"
	@echo "  make dev             - Run in development mode with hot reload"
	@echo "  make test            - Run all tests"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo "  make lint            - Lint the code"
	@echo "  make vuln-check      - Check for security vulnerabilities"
	@echo "  make swagger         - Generate Swagger documentation"
	@echo "  make mock            - Generate mocks for testing"
	@echo "  make migration-up    - Run database migrations up"
	@echo "  make migration-down  - Run database migrations down"
	@echo "  make install-hooks   - Install git hooks"
	@echo "  make encrypt-config  - Encrypt configuration value"
	@echo "  make test-config-reload - Test configuration hot reload"
	@echo "  make setup-encrypted-configs - Set up encrypted configuration values for all environments"
	@echo "  make help            - Show this help" 