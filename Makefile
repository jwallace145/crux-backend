# ===========================================
# Makefile for Crux Project Local Development
# ===========================================

# =========
# Variables
# =========

# Application Details
APP := crux_project.go
PROJECT_NAME := crux-project

# Docker Settings
COMPOSE_FILE := docker-compose.yml
DOCKERFILE_DEV := Dockerfile.dev


# Crux Project DB Settings
DB_CONTAINER := crux-project-db
DB_USER := crux_project_db_admin
DB_NAME := crux_project_db

# AWS Settings
AWS_REGION := us-east-1
ECR_REGISTRY := 650503560686.dkr.ecr.us-east-1.amazonaws.com
ECR_REPOSITORY := crux-api
ECS_CLUSTER := crux-api-cluster-dev
ECS_SERVICE := crux-api-service-dev
ECS_SERVICE_LOG_GROUP := crux-api-logs-dev

# Go files for formatting and linting
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")


# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_RED := \033[31m

.PHONY: start stop status restart logs
.PHONY: db-logs db-shell db-wait db-migrate db-reset
.PHONY: api-shell api-logs api-build test lint fmt fmt-check vet pre-commit
.PHONY: tf-fmt tf-fmt-check tf-plan tf-apply tf-destroy tf-output
.PHONY: repo-login image-build image-push service-deploy service-status service-logs deploy
.PHONY: check-deps

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "$(COLOR_BOLD)$(PROJECT_NAME) - Command Shortcuts$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_GREEN)Application:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make start" "Start application containers (API & DB)"
	@printf "  %-20s - %s\n" "make stop" "Stop application containers (API & DB)"
	@printf "  %-20s - %s\n" "make status" "Show status of application containers (API & DB)"
	@printf "  %-20s - %s\n" "make restart" "Restart application containers (API & DB)"
	@printf "  %-20s - %s\n" "make logs" "Show application containers logs (API & DB)"
	@echo ""
	@echo "$(COLOR_GREEN)Database:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make db-logs" "Show database container logs"
	@printf "  %-20s - %s\n" "make db-shell" "Start database container shell"
	@printf "  %-20s - %s\n" "make db-wait" "Wait for database container to start"
	@printf "  %-20s - %s\n" "make db-migrate" "Migrate database container tables"
	@printf "  %-20s - %s\n" "make db-reset" "(DANGER) Reset database and drop all data"
	@echo ""
	@echo "$(COLOR_GREEN)API:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make api-shell" "Open shell in API container"
	@printf "  %-20s - %s\n" "make api-logs" "Show API container logs"
	@printf "  %-20s - %s\n" "make api-build" "Build the API application binary"
	@printf "  %-20s - %s\n" "make test" "Run API tests"
	@printf "  %-20s - %s\n" "make lint" "Run linter (requires golangci-lint)"
	@printf "  %-20s - %s\n" "make fmt" "Format Go code"
	@printf "  %-20s - %s\n" "make fmt-check" "Check code formatting"
	@printf "  %-20s - %s\n" "make vet" "Run go vet"
	@printf "  %-20s - %s\n" "make pre-commit" "Run all pre-commit checks"
	@echo ""
	@echo "$(COLOR_GREEN)Terraform:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make tf-fmt" "Format Terraform code"
	@printf "  %-20s - %s\n" "make tf-fmt-check" "Check Terraform formatting"
	@printf "  %-20s - %s\n" "make tf-plan" "Plan the infrastructure changes before applying"
	@printf "  %-20s - %s\n" "make tf-apply" "Apply the infrastructure changes"
	@printf "  %-20s - %s\n" "make tf-destroy" "Destroy infrastructure"
	@printf "  %-20s - %s\n" "make tf-output" "Show Terraform outputs"
	@echo ""
	@echo "$(COLOR_GREEN)CICD:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make repo-login" "Authenticate with API image repository"
	@printf "  %-20s - %s\n" "make image-build" "Build API image for repository"
	@printf "  %-20s - %s\n" "make image-push" "Build and Push API image to repository"
	@printf "  %-20s - %s\n" "make service-deploy" "Start new API service deployment"
	@printf "  %-20s - %s\n" "make service-status" "Get latest API service status"
	@printf "  %-20s - %s\n" "make service-logs" "Get the latest API service logs"
	@printf "  %-20s - %s\n" "make deploy" "Build and push API image to repository and deploy updated service"
	@echo ""
	@echo "$(COLOR_GREEN)Utilities:$(COLOR_RESET)"
	@printf "  %-20s - %s\n" "make check-deps" "Check for required dependencies"

check-deps:
	@echo "$(COLOR_BOLD)Checking dependencies...$(COLOR_RESET)"
	@command -v docker >/dev/null 2>&1 || { echo "$(COLOR_RED)Error: docker is not installed$(COLOR_RESET)"; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "$(COLOR_RED)Error: docker-compose is not installed$(COLOR_RESET)"; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo "$(COLOR_RED)Error: go is not installed$(COLOR_RESET)"; exit 1; }
	@echo "$(COLOR_GREEN)✓ All dependencies found$(COLOR_RESET)"

start: check-deps
	@echo "$(COLOR_BOLD)Starting all containers...$(COLOR_RESET)"
	@TAG=$$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0") && \
	BRANCH=$$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown") && \
	APP_VERSION="$$TAG-$$BRANCH" && \
	echo "$(COLOR_BLUE)Using APP_VERSION: $$APP_VERSION$(COLOR_RESET)" && \
	APP_VERSION=$$APP_VERSION docker-compose -f $(COMPOSE_FILE) up -d --build
	@echo "$(COLOR_GREEN)✓ All containers started$(COLOR_RESET)"
	@$(MAKE) db-wait

stop:
	@echo "$(COLOR_BOLD)Stopping all containers...$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) down
	@echo "$(COLOR_GREEN)✓ All containers stopped$(COLOR_RESET)"

status:
	@echo "$(COLOR_BOLD)Container Status:$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) ps

restart:
	@echo "$(COLOR_BOLD)Restarting all containers...$(COLOR_RESET)"
	@$(MAKE) stop
	@$(MAKE) start

logs:
	@docker-compose -f $(COMPOSE_FILE) logs -f

db-logs:
	@echo "$(COLOR_BOLD)Showing database logs...$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) logs -f db

db-shell:
	@echo "$(COLOR_BOLD)Opening PostgreSQL shell...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Tip: Use \\dt to list tables, \\q to quit$(COLOR_RESET)"
	@docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

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

db-migrate: start db-wait
	@echo "$(COLOR_BOLD)Migrating database schema...$(COLOR_RESET)"
	@go run $(APP) &
	@sleep 3
	@pkill -f "go run $(APP)" || true
	@echo "$(COLOR_GREEN)✓ Database schema initialized$(COLOR_RESET)"

db-reset:
	@echo "$(COLOR_RED)$(COLOR_BOLD)WARNING: This will DELETE ALL DATA in the database!$(COLOR_RESET)"
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	@echo "$(COLOR_BOLD)Resetting database...$(COLOR_RESET)"
	@docker exec $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@echo "$(COLOR_GREEN)✓ Database reset complete$(COLOR_RESET)"
	@$(MAKE) db-migrate

api-shell:
	@echo "$(COLOR_BOLD)Opening shell in API container...$(COLOR_RESET)"
	@docker exec -it crux-project-api /bin/sh

api-logs:
	@echo "$(COLOR_BOLD)Showing API logs...$(COLOR_RESET)"
	@docker-compose -f $(COMPOSE_FILE) logs -f api

api-build: check-deps
	@echo "$(COLOR_BOLD)Building application...$(COLOR_RESET)"
	@go build -o bin/$(PROJECT_NAME) $(APP)
	@echo "$(COLOR_GREEN)✓ Build complete: bin/$(PROJECT_NAME)$(COLOR_RESET)"

test: check-deps
	@echo "$(COLOR_BOLD)Running tests...$(COLOR_RESET)"
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "$(COLOR_GREEN)✓ Tests complete$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)Coverage report:$(COLOR_RESET)"
	@go tool cover -func=coverage.out | tail -1

lint:
	@echo "$(COLOR_BOLD)Running CruxBackend Go linter...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(COLOR_GREEN)✓ Linting complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)Warning: golangci-lint not installed, skipping...$(COLOR_RESET)"; \
		echo "Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

fmt:
	@echo "$(COLOR_BOLD)Formatting CruxBackend Go source code and sorting imports...$(COLOR_RESET)"
	@gofmt -s -w $(GO_FILES)
	@goimports -local github.com/jwallace145/crux-backend -w $(GO_FILES)
	@echo "$(COLOR_GREEN)✓ Formatting complete$(COLOR_RESET)"

fmt-check:
	@echo "$(COLOR_BOLD)Checking CruxBackend Go source code formatting...$(COLOR_RESET)"
	@if [ -n "$$(gofmt -l $(GO_FILES))" ]; then \
		echo "$(COLOR_RED)✗ Code is not formatted. Run 'make fmt'$(COLOR_RESET)"; \
		gofmt -l $(GO_FILES); \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)✓ Code is properly formatted$(COLOR_RESET)"

vet:
	@echo "$(COLOR_BOLD)Running go vet...$(COLOR_RESET)"
	@go vet ./...
	@echo "$(COLOR_GREEN)✓ Vet complete$(COLOR_RESET)"

pre-commit: fmt-check vet lint test
	@echo "$(COLOR_GREEN)$(COLOR_BOLD)✓ All pre-commit checks passed!$(COLOR_RESET)"

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

tf-plan:
	@./infra/scripts/deploy.sh plan dev

tf-apply:
	@./infra/scripts/deploy.sh apply dev

tf-destroy:
	@./infra/scripts/deploy.sh destroy dev

tf-output:
	@./infra/scripts/deploy.sh output dev

repo-login:
	@echo "$(COLOR_BOLD)Authenticating with AWS ECR...$(COLOR_RESET)"
	@if ! command -v aws >/dev/null 2>&1; then \
		echo "$(COLOR_RED)Error: AWS CLI is not installed$(COLOR_RESET)"; \
		echo "Install from: https://aws.amazon.com/cli/"; \
		exit 1; \
	fi
	@aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR_REGISTRY)
	@echo "$(COLOR_GREEN)✓ ECR authentication successful$(COLOR_RESET)"

image-build: check-deps
	@echo "$(COLOR_BOLD)Building Docker image from $(DOCKERFILE_DEV) for linux/amd64...$(COLOR_RESET)"
	@docker build --platform linux/amd64 -f $(DOCKERFILE_DEV) -t $(ECR_REPOSITORY) .
	@docker tag $(ECR_REPOSITORY):latest $(ECR_REGISTRY)/$(ECR_REPOSITORY):latest
	@echo "$(COLOR_GREEN)✓ Docker image built and tagged for AWS (linux/amd64)$(COLOR_RESET)"

image-push: repo-login image-build
	@echo "$(COLOR_BOLD)Pushing image to ECR...$(COLOR_RESET)"
	@docker push $(ECR_REGISTRY)/$(ECR_REPOSITORY):latest
	@echo "$(COLOR_GREEN)✓ Image pushed to ECR: $(ECR_REGISTRY)/$(ECR_REPOSITORY):latest$(COLOR_RESET)"

service-deploy:
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

service-status:
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

service-logs:
	@echo "$(COLOR_BOLD)Fetching ECS task logs...$(COLOR_RESET)"
	@if ! command -v aws >/dev/null 2>&1; then \
		echo "$(COLOR_RED)Error: AWS CLI is not installed$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_YELLOW)Log group: /ecs/$(ECS_SERVICE_LOG_GROUP)$(COLOR_RESET)"
	@aws logs tail /ecs/$(ECS_SERVICE_LOG_GROUP) \
		--follow \
		--region $(AWS_REGION) \
		--format short

deploy: image-push service-deploy
	@echo ""
	@echo "$(COLOR_GREEN)$(COLOR_BOLD)✓ Complete deployment pipeline finished!$(COLOR_RESET)"
