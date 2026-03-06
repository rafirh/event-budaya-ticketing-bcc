.PHONY: build run dev test clean migrate seed fresh docker-up docker-down docker-build

# Application
APP_NAME=event-budaya-ticketing
MAIN_FILE=./cmd/main.go

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) $(MAIN_FILE)

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@./bin/$(APP_NAME)

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Running in development mode with hot reload..."
	@air

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf tmp/
	@rm -f coverage.out coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Run database migrations
migrate:
	@echo "Running migrations..."
	@go run $(MAIN_FILE) -migrate

# Run database seeders
seed:
	@echo "Running seeders..."
	@go run $(MAIN_FILE) -seed

# Fresh migration (drop all tables and re-run migrations)
fresh:
	@echo "Running fresh migration..."
	@go run $(MAIN_FILE) -fresh

# Run migrations and seeders
migrate-seed: migrate seed

# Fresh migration with seeders
fresh-seed: fresh seed

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker-compose build

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	@docker-compose logs -f

docker-restart: docker-down docker-up

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run

# Generate swagger docs (requires swag: go install github.com/swaggo/swag/cmd/swag@latest)
swagger:
	@echo "Generating swagger docs..."
	@swag init -g cmd/main.go -o docs

# Help
help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Build and run the application"
	@echo "  make dev           - Run with hot reload (requires air)"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make deps          - Download dependencies"
	@echo "  make migrate       - Run database migrations"
	@echo "  make seed          - Run database seeders"
	@echo "  make fresh         - Fresh migration (drop and re-run)"
	@echo "  make fresh-seed    - Fresh migration with seeders"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make docker-logs   - Show Docker logs"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Lint code"
	@echo "  make swagger       - Generate swagger docs"
