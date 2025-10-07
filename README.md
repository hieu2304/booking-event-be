# Event Booking API

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd event-booking-be
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Install development tools** (optional)
   ```bash
   make install-tools
   ```

4. **Setup environment**
   ```bash
   cp example.env .env
   # Edit .env with your configuration
   ```

5. **Start database services**
   ```bash
   make docker-up
   ```

6. **Run the application**
   ```bash
   make run
   # OR with hot reload
   make dev
   ```

The API will be available at `http://localhost:8080`

## Makefile Commands

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Run the application
make dev               # Run with hot reload (requires air)
make test              # Run tests
make test-coverage     # Run tests with coverage report
make lint              # Run linter
make clean             # Clean build artifacts
make docker-up         # Start Docker containers
make docker-down       # Stop Docker containers
make docker-logs       # View Docker logs
make migrate-up        # Run database migrations up
make migrate-down      # Run database migrations down
make migrate-create    # Create new migration
make deps              # Download dependencies
make tidy              # Tidy go modules
make install-tools     # Install development tools
make setup             # Full setup (deps + tools)
```

## Development

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage
```

### Hot Reload

```bash
# Install air
make install-tools

# Run with hot reload
make dev
```

### Database Migrations

```bash
# Create new migration
make migrate-create

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down
```

## Docker

### Start Services

```bash
docker-compose up -d
```

### Stop Services

```bash
docker-compose down
```

### View Logs

```bash
make docker-logs
```


## Architecture

This project follows Clean Architecture principles with clear separation of concerns:

- **Models Layer**: GORM entities and DTOs
- **Repository Layer**: Database operations (CRUD)
- **Service Layer**: Business logic and transactions
- **Handler Layer**: HTTP request/response handling
- **Middleware Layer**: Authentication, logging, error handling
