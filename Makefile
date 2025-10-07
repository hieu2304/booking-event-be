.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down lint

APP_NAME=event-booking-api
BUILD_DIR=bin
MAIN_PATH=cmd/api/main.go

help:
	@echo "Available commands:"
	@echo "  make build           - Build the application"
	@echo "  make run             - Run the application"
	@echo "  make dev             - Run with hot reload (requires air)"
	@echo "  make test            - Run tests"
	@echo "  make test-coverage   - Run tests with coverage"
	@echo "  make lint            - Run linter"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make docker-up       - Start Docker containers"
	@echo "  make docker-down     - Stop Docker containers"
	@echo "  make migrate-up      - Run database migrations up"
	@echo "  make migrate-down    - Run database migrations down"
	@echo "  make deps            - Download dependencies"
	@echo "  make tidy            - Tidy go modules"

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)

dev:
	@echo "Running with hot reload..."
	@air

test:
	@echo "Running tests..."
	@go test -v ./tests/...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./tests/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run specific test
test-run:
	@go test -v ./tests/... -run $(TEST)

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Containers started"

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Containers stopped"

docker-logs:
	@docker-compose logs -f

migrate-up:
	@echo "Running migrations up..."
	@migrate -path migrations -database "${DATABASE_URL}" up

migrate-down:
	@echo "Running migrations down..."
	@migrate -path migrations -database "${DATABASE_URL}" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

deps:
	@echo "Downloading dependencies..."
	@go mod download

tidy:
	@echo "Tidying go modules..."
	@go mod tidy

install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed"

setup: deps install-tools
	@echo "Setup complete"

.DEFAULT_GOAL := help

# Docker - Build app image
docker-build:
	@echo "Building Docker image..."
	@docker build -t event-booking-api:latest .
	@echo "✅ Image built: event-booking-api:latest"

docker-build-prod:
	@echo "Building production Docker image..."
	@docker build --no-cache -t event-booking-api:prod .
	@echo "✅ Production image built"

docker-run-app:
	@echo "Running app in Docker..."
	@docker run --rm -p 8080:8080 \
		--env-file .env \
		--name event-booking-api \
		event-booking-api:latest

docker-stop-app:
	@docker stop event-booking-api || true

# Run everything with docker-compose (including app)
docker-up-all:
	@echo "Starting all services..."
	@docker-compose up -d --build
	@echo "✅ All services started"

docker-app-logs:
	@docker-compose logs -f app

# CI/CD helpers
ci-test:
	@echo "Running CI tests..."
	@go test -v -race -coverprofile=coverage.out ./...

ci-lint:
	@echo "Running CI lint..."
	@golangci-lint run --timeout=5m ./...

ci-build:
	@echo "Building for CI..."
	@CGO_ENABLED=0 go build -o bin/$(APP_NAME) $(MAIN_PATH)
