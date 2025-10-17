# Makefile for Crux Backend Local Development
# ===========================================

# Variables
COMPOSE_FILE := docker-compose.yml
APP := main.go
PROJECT_NAME := crux-backend
DB_CONTAINER := cruxdb
DB_USER := crux_user
DB_NAME := cruxdb
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")

# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_RED := \033[31m

.PHONY: help up down restart logs logs-api logs-db status clean build test lint fmt fmt-check vet pre-commit run bootstrap reset db-wait db-shell db-migrate db-reset-force check-deps test-api test-db test-all api-shell test-api test-db test-all api-shell

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "$(COLOR_BOLD)$(PROJECT_NAME) - Available Commands$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_GREEN)Database:$(COLOR_RESET)"
	@echo "  make up              - Start PostgreSQL container"
	@echo "  make down            - Stop PostgreSQL container"
	@echo "  make restart         - Restart PostgreSQL container"
	@echo "  make logs            - Show logs from all containers"
	@echo "  make logs-api        - Show API container logs only"
	@echo "  make logs-db         - Show database container logs only"
	@echo "  make status          - Show container status"
	@echo "  make db-wait         - Wait for database to be ready"
	@echo "  make db-shell        - Open PostgreSQL shell"
	@echo ""
	@echo "$(COLOR_GREEN)Application:$(COLOR_RESET)"
	@echo "  make run             - Run Fiber API (starts DB if needed)"
	@echo "  make build           - Build the application binary"
	@echo "  make test            - Run tests"
	@echo "  make lint            - Run linter (requires golangci-lint)"
	@echo "  make fmt             - Format Go code"
	@echo "  make fmt-check       - Check code formatting"
	@echo "  make vet             - Run go vet"
	@echo "  make pre-commit      - Run all pre-commit checks"
	@echo "  make bootstrap       - Initialize database schema"
	@echo ""
	@echo "$(COLOR_GREEN)Database Management:$(COLOR_RESET)"
	@echo "  make db-migrate      - Run migrations (same as bootstrap)"
	@echo "  make reset           - Reset database (DESTRUCTIVE - prompts for confirmation)"
	@echo "  make db-reset-force  - Reset database without confirmation (DANGEROUS)"
	@echo ""
	@echo "$(COLOR_GREEN)Cleanup:$(COLOR_RESET)"
	@echo "  make clean           - Remove binaries and temporary files"
	@echo "  make clean-all       - Remove containers, volumes, and binaries"
	@echo ""
	@echo "$(COLOR_GREEN)Utilities:$(COLOR_RESET)"
	@echo "  make check-deps      - Check for required dependencies"
	@echo ""
	@echo "$(COLOR_GREEN)Testing:$(COLOR_RESET)"
	@echo "  make test-api        - Test API health endpoint"
	@echo "  make test-db         - Test database connection"
	@echo "  make test-all        - Test both API and database"

## check-deps: Verify required tools are installed
check-deps:
	@echo "$(COLOR_BOLD)Checking dependencies...$(COLOR_RESET)"
	@command -v docker >/dev/null 2>&1 || { echo "$(COLOR_RED)Error: docker is not installed$(COLOR_RESET)"; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "$(COLOR_RED)Error: docker-compose is not installed$(COLOR_RESET)"; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo "$(COLOR_RED)Error: go is not installed$(COLOR_RESET)"; exit 1; }
	@echo "$(COLOR_GREEN)✓ All dependencies found$(COLOR_RESET)"

## up: Start all containers (database and API)
up: check-deps
	@echo "$(COLOR_BOLD)Starting all containers...$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) up -d --build
	@echo "$(COLOR_GREEN)✓ All containers started$(COLOR_RESET)"
	@$(MAKE) db-wait

## down: Stop all containers
down:
	@echo "$(COLOR_BOLD)Stopping all containers...$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) down
	@echo "$(COLOR_GREEN)✓ All containers stopped$(COLOR_RESET)"

## restart: Restart all containers
restart:
	@echo "$(COLOR_BOLD)Restarting all containers...$(COLOR_RESET)"
	@$(MAKE) down
	@$(MAKE) up

## logs: Show container logs (both API and database)
logs:
	@docker-compose -f $(COMPOSE_FILE) logs -f

## logs-api: Show API container logs only
logs-api:
	@echo "$(COLOR_BOLD)Showing API logs...$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) logs -f api

## logs-db: Show database container logs only
logs-db:
	@echo "$(COLOR_BOLD)Showing database logs...$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) logs -f db

## status: Show container status
status:
	@echo "$(COLOR_BOLD)Container Status:$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) ps

## db-wait: Wait for database to be ready
db-wait:
	@echo "$(COLOR_BOLD)Waiting for database to be ready...$(COLOR_RESET)"
	@timeout=60; \
	elapsed=0; \
	while ! docker exec $(DB_CONTAINER) pg_isready -U $(DB_USER) -d $(DB_NAME) >/dev/null 2>&1; do \
		if [ $$elapsed -ge $$timeout ]; then \
			echo "$(COLOR_RED)✗ Database failed to start within $$timeout seconds$(COLOR_RESET)"; \
			exit 1; \
		fi; \
		printf "."; \
		sleep 1; \
		elapsed=$$((elapsed + 1)); \
	done; \
	echo ""; \
	echo "$(COLOR_GREEN)✓ Database is ready$(COLOR_RESET)"

## db-shell: Open PostgreSQL shell
db-shell:
	@echo "$(COLOR_BOLD)Opening PostgreSQL shell...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Tip: Use \\dt to list tables, \\q to quit$(COLOR_RESET)"
	@docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

## api-shell: Open shell in API container
api-shell:
	@echo "$(COLOR_BOLD)Opening shell in API container...$(COLOR_RESET)"
	@docker exec -it crux-api /bin/sh

## build: Build the application binary
build: check-deps
	@echo "$(COLOR_BOLD)Building application...$(COLOR_RESET)"
	@go build -o bin/$(PROJECT_NAME) $(APP)
	@echo "$(COLOR_GREEN)✓ Build complete: bin/$(PROJECT_NAME)$(COLOR_RESET)"

## test: Run tests
test: check-deps
	@echo "$(COLOR_BOLD)Running tests...$(COLOR_RESET)"
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "$(COLOR_GREEN)✓ Tests complete$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)Coverage report:$(COLOR_RESET)"
	@go tool cover -func=coverage.out | tail -1

## lint: Run linter
lint:
	@echo "$(COLOR_BOLD)Running CruxBackend Go linter...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(COLOR_GREEN)✓ Linting complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)Warning: golangci-lint not installed, skipping...$(COLOR_RESET)"; \
		echo "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## fmt: Format Go code and imports
fmt:
	@echo "$(COLOR_BOLD)Formatting CruxBackend Go source code and sorting imports...$(COLOR_RESET)"
	@gofmt -s -w $(GO_FILES)
	@goimports -local github.com/jwallace145/crux-backend -w $(GO_FILES)
	@echo "$(COLOR_GREEN)✓ Formatting complete$(COLOR_RESET)"

## fmt-check: Check if code is formatted
fmt-check:
	@echo "$(COLOR_BOLD)Checking CruxBackend Go source code formatting...$(COLOR_RESET)"
	@if [ -n "$$(gofmt -l $(GO_FILES))" ]; then \
		echo "$(COLOR_RED)✗ Code is not formatted. Run 'make fmt'$(COLOR_RESET)"; \
		gofmt -l $(GO_FILES); \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)✓ Code is properly formatted$(COLOR_RESET)"

## vet: Run go vet
vet:
	@echo "$(COLOR_BOLD)Running go vet...$(COLOR_RESET)"
	@go vet ./...
	@echo "$(COLOR_GREEN)✓ Vet complete$(COLOR_RESET)"

## pre-commit: Run all pre-commit checks locally
pre-commit: fmt-check vet lint test
	@echo "$(COLOR_GREEN)$(COLOR_BOLD)✓ All pre-commit checks passed!$(COLOR_RESET)"

## run: Run Fiber API (ensures database is running)
run: up
	@echo "$(COLOR_BOLD)Starting Fiber API...$(COLOR_RESET)"
	@go run $(APP)

## bootstrap: Initialize database schema
bootstrap: up db-wait
	@echo "$(COLOR_BOLD)Bootstrapping database schema...$(COLOR_RESET)"
	@go run $(APP) &
	@sleep 3
	@pkill -f "go run $(APP)" || true
	@echo "$(COLOR_GREEN)✓ Database schema initialized$(COLOR_RESET)"

## db-migrate: Run migrations (alias for bootstrap)
db-migrate: bootstrap

## reset: Reset database with confirmation (DESTRUCTIVE)
reset:
	@echo "$(COLOR_RED)$(COLOR_BOLD)WARNING: This will DELETE ALL DATA in the database!$(COLOR_RESET)"
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	@$(MAKE) db-reset-force

## db-reset-force: Reset database without confirmation (DANGEROUS)
db-reset-force: up db-wait
	@echo "$(COLOR_BOLD)Resetting database...$(COLOR_RESET)"
	@docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@echo "$(COLOR_GREEN)✓ Database schema dropped$(COLOR_RESET)"
	@$(MAKE) bootstrap

## clean: Remove binaries and temporary files
clean:
	@echo "$(COLOR_BOLD)Cleaning temporary files...$(COLOR_RESET)"
	@rm -rf bin/
	@rm -f coverage.out
	@go clean
	@echo "$(COLOR_GREEN)✓ Cleanup complete$(COLOR_RESET)"

## clean-all: Remove containers, volumes, and binaries
clean-all: clean down
	@echo "$(COLOR_BOLD)Removing Docker volumes...$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) down -v
	@echo "$(COLOR_GREEN)✓ All resources removed$(COLOR_RESET)"

## test-api: Test API health endpoint
test-api:
	@echo "$(COLOR_BOLD)Testing API health endpoint...$(COLOR_RESET)"
	@if curl -sf http://localhost:3000/health > /dev/null 2>&1; then \
		echo "$(COLOR_GREEN)✓ API is responding$(COLOR_RESET)"; \
		echo ""; \
		echo "$(COLOR_BOLD)Full response:$(COLOR_RESET)"; \
		curl -s http://localhost:3000/health; \
		echo ""; \
	else \
		echo "$(COLOR_RED)✗ API is not responding$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Hint: Run 'make up' to start the API$(COLOR_RESET)"; \
		exit 1; \
	fi

## test-db: Test database connection
test-db:
	@echo "$(COLOR_BOLD)Testing database connection...$(COLOR_RESET)"
	@if docker exec $(DB_CONTAINER) pg_isready -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; then \
		echo "$(COLOR_GREEN)✓ Database is responding$(COLOR_RESET)"; \
		echo ""; \
		echo "$(COLOR_BOLD)Database info:$(COLOR_RESET)"; \
		docker exec $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT version();" 2>/dev/null | head -3; \
		echo ""; \
		echo "$(COLOR_BOLD)Tables:$(COLOR_RESET)"; \
		docker exec $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "\\dt" 2>/dev/null; \
	else \
		echo "$(COLOR_RED)✗ Database is not responding$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Hint: Run 'make up' to start the database$(COLOR_RESET)"; \
		exit 1; \
	fi

## test-all: Test both API and database
test-all:
	@echo "$(COLOR_BOLD)=== Running All Tests ===$(COLOR_RESET)"
	@echo ""
	@$(MAKE) test-db
	@echo ""
	@$(MAKE) test-api
	@echo ""
	@echo "$(COLOR_GREEN)$(COLOR_BOLD)✓ All systems operational!$(COLOR_RESET)"
