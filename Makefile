# ===========================================
# Makefile for Crux Backend Local Development
# ===========================================

# Variables
COMPOSE_FILE := docker-compose.yml
APP := main.go
PROJECT_NAME := crux-backend
DB_CONTAINER := cruxdb
DB_USER := crux_user
DB_NAME := cruxdb
AWS_REGION := us-east-1
ECR_REGISTRY := 650503560686.dkr.ecr.us-east-1.amazonaws.com
ECR_REPOSITORY := crux-api
DOCKERFILE_DEV := Dockerfile.dev
ECS_CLUSTER := crux-api-cluster-dev
ECS_SERVICE := crux-api-service-dev
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")


# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_RED := \033[31m

.PHONY: help up down restart logs logs-api logs-db status clean build test lint fmt fmt-check tf-fmt tf-fmt-check tf-plan tf-apply vet pre-commit run bootstrap reset db-wait db-shell db-migrate db-reset-force ecr-login ecr-build ecr-push ecr-deploy ecs-deploy ecs-status ecs-logs deploy check-deps test-api test-db test-all api-shell test-api test-db test-all api-shell

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "$(COLOR_BOLD)$(PROJECT_NAME) - Available Commands$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_GREEN)Database:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make up" "Start PostgreSQL container"
	@printf "  %-20s - %s\n" "make down" "Stop PostgreSQL container"
	@printf "  %-20s - %s\n" "make restart" "Restart PostgreSQL container"
	@printf "  %-20s - %s\n" "make logs" "Show logs from all containers"
	@printf "  %-20s - %s\n" "make logs-api" "Show API container logs only"
	@printf "  %-20s - %s\n" "make logs-db" "Show database container logs only"
	@printf "  %-20s - %s\n" "make status" "Show container status"
	@printf "  %-20s - %s\n" "make db-wait" "Wait for database to be ready"
	@printf "  %-20s - %s\n" "make db-shell" "Open PostgreSQL shell"
	@echo ""
	@echo "$(COLOR_GREEN)Application:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make run" "Run Fiber API (starts DB if needed)"
	@printf "  %-20s - %s\n" "make build" "Build the application binary"
	@printf "  %-20s - %s\n" "make test" "Run tests"
	@printf "  %-20s - %s\n" "make lint" "Run linter (requires golangci-lint)"
	@printf "  %-20s - %s\n" "make fmt" "Format Go code"
	@printf "  %-20s - %s\n" "make fmt-check" "Check code formatting"
	@printf "  %-20s - %s\n" "make vet" "Run go vet"
	@printf "  %-20s - %s\n" "make pre-commit" "Run all pre-commit checks"
	@printf "  %-20s - %s\n" "make bootstrap" "Initialize database schema"
	@echo ""
	@echo "$(COLOR_GREEN)Terraform:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make tf-fmt" "Format Terraform code"
	@printf "  %-20s - %s\n" "make tf-fmt-check" "Check Terraform formatting"
	@printf "  %-20s - %s\n" "make tf-plan" "Plan the infrastructure changes before applying"
	@printf "  %-20s - %s\n" "make tf-apply" "Apply the infrastructure changes"
	@printf "  %-20s - %s\n" "make tf-destroy" "Destroy infrastructure"
	@printf "  %-20s - %s\n" "make tf-output" "Show Terraform outputs"
	@echo ""
	@echo "$(COLOR_GREEN)Database Management:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make db-shell" "Open PostgreSQL interactive shell"
	@printf "  %-20s - %s\n" "make db-migrate" "Run migrations (same as bootstrap)"
	@printf "  %-20s - %s\n" "make reset" "Reset database (DESTRUCTIVE - prompts for confirmation)"
	@printf "  %-20s - %s\n" "make db-reset-force" "Reset database without confirmation (DANGEROUS)"
	@echo ""
	@echo "$(COLOR_GREEN)ECR Commands:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make ecr-login" "Authenticate with AWS ECR"
	@printf "  %-20s - %s\n" "make ecr-build" "Build Docker image for ECR"
	@printf "  %-20s - %s\n" "make ecr-push" "Push Docker image to ECR"
	@printf "  %-20s - %s\n" "make ecr-deploy" "Full ECR deployment pipeline"
	@echo ""
	@echo "$(COLOR_GREEN)Cleanup:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make clean" "Remove binaries and temporary files"
	@printf "  %-20s - %s\n" "make clean-all" "Remove containers, volumes, and binaries"
	@echo ""
	@echo "$(COLOR_GREEN)Utilities:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make check-deps" "Check for required dependencies"
	@echo ""
	@echo "$(COLOR_GREEN)Testing:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make test-api" "Test API health endpoint"
	@printf "  %-20s - %s\n" "make test-db" "Test database connection"
	@printf "  %-20s - %s\n" "make test-all" "Test both API and database"
	@printf "  %-20s - %s\n" "make api-shell" "Open shell in API container"m

## check-deps: Verify required tools are installed
check-deps:
	@echo "$(COLOR_BOLD)Checking dependencies...$(COLOR_RESET)"
	@command -v docker >/dev/null 2>&1 || { echo "$(COLOR_RED)Error: docker is not installed$(COLOR_RESET)"; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "$(COLOR_RED)Error: docker-compose is not installed$(COLOR_RESET)"; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo "$(COLOR_RED)Error: go is not installed$(COLOR_RESET)"; exit 1; }
	@echo "$(COLOR_GREEN)✓ All dependencies found$(COLOR_RESET)"

## up: Start all containers (db and API)
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

## logs: Show container logs (both API and db)
logs:
	@docker-compose -f $(COMPOSE_FILE) logs -f

## logs-api: Show API container logs only
logs-api:
	@echo "$(COLOR_BOLD)Showing API logs...$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) logs -f api

## logs-db: Show db container logs only
logs-db:
	@echo "$(COLOR_BOLD)Showing database logs...$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) logs -f db

## status: Show container status
status:
	@echo "$(COLOR_BOLD)Container Status:$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) ps

## db-wait: Wait for db to be ready
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

## ecr-login: Authenticate with AWS ECR
ecr-login:
	@echo "$(COLOR_BOLD)Authenticating with AWS ECR...$(COLOR_RESET)"
	@if ! command -v aws >/dev/null 2>&1; then \
		echo "$(COLOR_RED)Error: AWS CLI is not installed$(COLOR_RESET)"; \
		echo "Install from: https://aws.amazon.com/cli/"; \
		exit 1; \
	fi
	@aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR_REGISTRY)
	@echo "$(COLOR_GREEN)✓ ECR authentication successful$(COLOR_RESET)"

## ecr-build: Build Docker image for ECR (multi-platform for AWS compatibility)
ecr-build: check-deps
	@echo "$(COLOR_BOLD)Building Docker image from $(DOCKERFILE_DEV) for linux/amd64...$(COLOR_RESET)"
	@docker build --platform linux/amd64 -f $(DOCKERFILE_DEV) -t $(ECR_REPOSITORY) .
	@docker tag $(ECR_REPOSITORY):latest $(ECR_REGISTRY)/$(ECR_REPOSITORY):latest
	@echo "$(COLOR_GREEN)✓ Docker image built and tagged for AWS (linux/amd64)$(COLOR_RESET)"

## ecr-push: Push Docker image to ECR (includes login and build)
ecr-push: ecr-login ecr-build
	@echo "$(COLOR_BOLD)Pushing image to ECR...$(COLOR_RESET)"
	@docker push $(ECR_REGISTRY)/$(ECR_REPOSITORY):latest
	@echo "$(COLOR_GREEN)✓ Image pushed to ECR: $(ECR_REGISTRY)/$(ECR_REPOSITORY):latest$(COLOR_RESET)"

## ecr-deploy: Full ECR deployment pipeline (login, build, push)
ecr-deploy: ecr-push
	@echo "$(COLOR_GREEN)$(COLOR_BOLD)✓ ECR deployment complete!$(COLOR_RESET)"

## ecs-deploy: Force ECS service to deploy new image
ecs-deploy:
	@echo "$(COLOR_BOLD)Forcing ECS service deployment...$(COLOR_RESET)"
	@if ! command -v aws >/dev/null 2>&1; then \
		echo "$(COLOR_RED)Error: AWS CLI is not installed$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_YELLOW)Updating service: $(ECS_SERVICE) in cluster: $(ECS_CLUSTER)$(COLOR_RESET)"
	@aws ecs update-service \
		--cluster $(ECS_CLUSTER) \
		--service $(ECS_SERVICE) \
		--force-new-deployment \
		--region $(AWS_REGION) \
		--no-cli-pager > /dev/null
	@echo "$(COLOR_GREEN)✓ ECS deployment initiated$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Tip: Run 'make ecs-status' to monitor deployment progress$(COLOR_RESET)"

## ecs-status: Check ECS service deployment status
ecs-status:
	@echo "$(COLOR_BOLD)Checking ECS service status...$(COLOR_RESET)"
	@if ! command -v aws >/dev/null 2>&1; then \
		echo "$(COLOR_RED)Error: AWS CLI is not installed$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo ""
	@echo "$(COLOR_BOLD)Service Status:$(COLOR_RESET)"
	@aws ecs describe-services \
		--cluster $(ECS_CLUSTER) \
		--services $(ECS_SERVICE) \
		--region $(AWS_REGION) \
		--query 'services[0].{Status:status,Running:runningCount,Desired:desiredCount,Pending:pendingCount}' \
		--output table
	@echo ""
	@echo "$(COLOR_BOLD)Recent Deployments:$(COLOR_RESET)"
	@aws ecs describe-services \
		--cluster $(ECS_CLUSTER) \
		--services $(ECS_SERVICE) \
		--region $(AWS_REGION) \
		--query 'services[0].deployments[*].{Status:status,Desired:desiredCount,Running:runningCount,Created:createdAt}' \
		--output table
	@echo ""
	@echo "$(COLOR_BOLD)Task Status:$(COLOR_RESET)"
	@aws ecs list-tasks \
		--cluster $(ECS_CLUSTER) \
		--service-name $(ECS_SERVICE) \
		--region $(AWS_REGION) \
		--query 'taskArns[*]' \
		--output text | xargs -I {} aws ecs describe-tasks \
		--cluster $(ECS_CLUSTER) \
		--tasks {} \
		--region $(AWS_REGION) \
		--query 'tasks[*].{TaskId:taskArn,Status:lastStatus,Health:healthStatus,Started:startedAt}' \
		--output table 2>/dev/null || echo "No running tasks"

## ecs-logs: Show recent ECS task logs
ecs-logs:
	@echo "$(COLOR_BOLD)Fetching ECS task logs...$(COLOR_RESET)"
	@if ! command -v aws >/dev/null 2>&1; then \
		echo "$(COLOR_RED)Error: AWS CLI is not installed$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_YELLOW)Log group: /ecs/$(ECS_SERVICE)$(COLOR_RESET)"
	@aws logs tail /ecs/$(ECS_SERVICE) \
		--follow \
		--region $(AWS_REGION) \
		--format short

## deploy: Complete deployment pipeline (build, push, deploy to ECS)
deploy: ecr-push ecs-deploy
	@echo ""
	@echo "$(COLOR_GREEN)$(COLOR_BOLD)✓ Complete deployment pipeline finished!$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_YELLOW)Next steps:$(COLOR_RESET)"
	@echo "  1. Run 'make ecs-status' to monitor deployment"
	@echo "  2. Run 'make ecs-logs' to watch application logs"
	@echo "  3. Test your API endpoint once deployment completes"

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

# tf-fmt: Format Terraform IaC
tf-fmt:
	@echo "$(COLOR_BOLD)Formatting Terraform code...$(COLOR_RESET)"
	@if command -v terraform >/dev/null 2>&1; then \
		terraform fmt -recursive infra/; \
		echo "$(COLOR_GREEN)✓ Terraform formatting complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_RED)Error: terraform is not installed$(COLOR_RESET)"; \
		echo "Install from: https://www.terraform.io/downloads"; \
		exit 1; \
	fi

## tf-fmt-check: Check Terraform formatting
tf-fmt-check:
	@echo "$(COLOR_BOLD)Checking Terraform code formatting...$(COLOR_RESET)"
	@if command -v terraform >/dev/null 2>&1; then \
		if terraform fmt -check -recursive infra/; then \
			echo "$(COLOR_GREEN)✓ Terraform code is properly formatted$(COLOR_RESET)"; \
		else \
			echo "$(COLOR_RED)✗ Terraform code is not formatted. Run 'make tf-fmt'$(COLOR_RESET)"; \
			exit 1; \
		fi; \
	else \
		echo "$(COLOR_RED)Error: terraform is not installed$(COLOR_RESET)"; \
		exit 1; \
	fi

## tf-plan: Plan Terraform changes
tf-plan:
	@./infra/scripts/deploy.sh plan dev

## tf-apply: Apply Terraform changes
tf-apply:
	@./infra/scripts/deploy.sh apply dev

## tf-destroy: Destroy infrastructure
tf-destroy:
	@./infra/scripts/deploy.sh destroy dev

## tf-output: Show Terraform outputs
tf-output:
	@./infra/scripts/deploy.sh output dev

## vet: Run go vet
vet:
	@echo "$(COLOR_BOLD)Running go vet...$(COLOR_RESET)"
	@go vet ./...
	@echo "$(COLOR_GREEN)✓ Vet complete$(COLOR_RESET)"

## pre-commit: Run all pre-commit checks locally
pre-commit: fmt-check vet lint test
	@echo "$(COLOR_GREEN)$(COLOR_BOLD)✓ All pre-commit checks passed!$(COLOR_RESET)"

## run: Run Fiber API (ensures db is running)
run: up
	@echo "$(COLOR_BOLD)Starting Fiber API...$(COLOR_RESET)"
	@go run $(APP)

## bootstrap: Initialize db schema
bootstrap: up db-wait
	@echo "$(COLOR_BOLD)Bootstrapping database schema...$(COLOR_RESET)"
	@go run $(APP) &
	@sleep 3
	@pkill -f "go run $(APP)" || true
	@echo "$(COLOR_GREEN)✓ Database schema initialized$(COLOR_RESET)"

## db-migrate: Run migrations (alias for bootstrap)
db-migrate: bootstrap

## reset: Reset db with confirmation (DESTRUCTIVE)
reset:
	@echo "$(COLOR_RED)$(COLOR_BOLD)WARNING: This will DELETE ALL DATA in the database!$(COLOR_RESET)"
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	@$(MAKE) db-reset-force

## db-reset-force: Reset db without confirmation (DANGEROUS)
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

## test-db: Test db connection
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

## test-all: Test both API and db
test-all:
	@echo "$(COLOR_BOLD)=== Running All Tests ===$(COLOR_RESET)"
	@echo ""
	@$(MAKE) test-db
	@echo ""
	@$(MAKE) test-api
	@echo ""
	@echo "$(COLOR_GREEN)$(COLOR_BOLD)✓ All systems operational!$(COLOR_RESET)"
