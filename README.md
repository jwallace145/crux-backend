# CruxBackend

[![Release and Deploy](https://github.com/jwallace145/crux-project/actions/workflows/release.yml/badge.svg?branch=main)](https://github.com/jwallace145/crux-project/actions/workflows/release.yml)

A high-performance rock climbing API built with Go and [Fiber](https://gofiber.io/) that powers web and mobile applications for the climbing community.

## Overview

CruxBackend is the central backend service for climbers to discover routes, log climbs, track progress, and connect with the climbing community. Built for scalability and deployed on AWS with containerized microservices.

## Core MVP Features

These are the essential features that define the product. Use this list to guide development priorities:

### 1. Natural Language Route Search
**Goal:** Enable climbers to find routes using conversational queries
**Example:** _"Give me all single pitch 5.8-5.9 trad routes within a 2 hour drive of Denver, CO"_
**Status:** 🚧 Not Started

**Key Requirements:**
- Parse natural language queries (difficulty, style, location, distance)
- Geocoding and distance calculations
- Filter by route attributes (pitch count, grade, type)
- Return ranked results

### 2. Climb Logging & Progress Tracking
**Goal:** Record climbs (outdoor & indoor) and track progress over time
**Features:**
- Log outdoor climbs with route association
- Log indoor gym sessions
- Track climbing partners
- View personal climbing history and statistics
- Filter climbs by date range

**Status:** ✅ Core logging implemented, analytics pending

### 3. Community & Social Features
**Goal:** Enable climbers to share experiences and find partners
**Features:**
- Comment on specific routes
- Comment on other climbers' logged climbs (if following)
- Follow/unfollow climbers
- Find climbing partners based on shared interests
- Route ratings and reviews

**Status:** 🚧 Not Started


## Architecture

### Local Development
```
┌─────────────────┐
│  Developer      │
└────────┬────────┘
         │
    ┌────▼─────────┐
    │   Fiber API  │ (Docker Container)
    │   Port 3000  │
    └────┬─────────┘
         │
    ┌────▼──────────┐
    │  PostgreSQL   │ (Docker Container)
    │   Port 5432   │
    └───────────────┘
```

### Production (AWS)
```
Internet
    │
    ▼
┌─────────────────────────┐
│   Route53 DNS           │
│   *-api.domain.com      │
└──────────┬──────────────┘
           │
    ┌──────▼──────────┐
    │  Application    │
    │  Load Balancer  │
    │  (ALB)          │
    └──────┬──────────┘
           │
    ┌──────▼────────────┐
    │  ECS Fargate      │
    │  - Go Fiber API   │
    │  - Auto-scaling   │
    │  - Health checks  │
    └──────┬────────────┘
           │
    ┌──────▼────────────┐
    │  RDS PostgreSQL   │
    │  (Private subnet) │
    └───────────────────┘

┌─────────────────────────┐
│  EC2 Bastion Host       │
│  (Network debugging)    │
└─────────────────────────┘
```

## Quick Start

### Prerequisites
- [Go 1.24+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.gnu.org/software/make/)

### Local Development

```bash
# Start all services (API + Database)
make up

# Or run API locally with auto-reload (starts DB container)
make run

# View logs
make logs          # All services
make logs-api      # API only
make logs-db       # Database only

# Run tests with coverage
make test

# Format code
make fmt

# Lint code
make lint

# Run all pre-commit checks
make pre-commit

# Database management
make db-shell      # Open PostgreSQL shell
make bootstrap     # Run migrations
make reset         # Reset db (prompts for confirmation)

# View all available commands
make help
```

The API will be available at `http://localhost:3000`

### Health Check
```bash
curl http://localhost:3000/health
```

### API Endpoints

#### Authentication
- `POST /login` - User login (returns JWT tokens)
- `POST /logout` - User logout (revokes session)
- `POST /refresh` - Refresh access token

#### Users
- `POST /users` - Create new user account

#### Climbs
- `POST /climbs` - Log a climb (outdoor or indoor)
- `GET /climbs?user_id=X&start_date=Y&end_date=Z` - Get user's climbs

## API Documentation

### Interactive API Documentation with Scalar

CruxBackend provides comprehensive, interactive API documentation powered by [Scalar](https://guides.scalar.com/scalar/introduction), designed for developers who want to explore, test, and integrate with the API.

**Live Documentation:** [https://jwallace145.github.io/crux-backend](https://jwallace145.github.io/crux-backend)

#### Features

The interactive documentation provides:
- **Live API Testing** - Test all endpoints directly from your browser
- **Multiple Environments** - Switch between local development (localhost:3000) and deployed servers (dev-api.cruxproject.io)
- **Request/Response Examples** - View example requests and responses for every endpoint
- **Schema Documentation** - Detailed documentation of all data models and validation rules
- **Authentication Flow** - Test the complete authentication workflow with JWT tokens
- **Search & Browse** - Quickly find endpoints organized by category (Health, Auth, Users, Climbs)
- **OpenAPI Specification** - Based on OpenAPI 3.0.3 standard for easy integration with other tools

#### Local Development

When running the API locally, the documentation is available at:
```
http://localhost:3000/docs
```

The OpenAPI specification file can be accessed at:
```
http://localhost:3000/docs/openapi.yaml
```

#### Using the Documentation

1. **Select Environment** - Use the server dropdown to choose between local development or deployed servers
2. **Browse Endpoints** - Navigate through all available API endpoints organized by tags
3. **Try It Out** - Click "Try It" on any endpoint to send real requests to the selected environment
4. **Authentication** - Login via `POST /login` to test authenticated endpoints (cookies are handled automatically)
5. **View Schemas** - Explore detailed request and response schemas with validation rules

**Environment Switching:**
The documentation supports testing against multiple environments without leaving the page. Simply select your target environment from the server dropdown:
- **Local Development** (`http://localhost:3000`) - Test changes during local development
- **Development Server** (`http://dev-api.cruxproject.io`) - Test deployed changes in the dev environment

#### For Developers

The OpenAPI specification can be used with various tools:
- Import into Postman, Insomnia, or other API clients
- Generate client SDKs using [OpenAPI Generator](https://openapi-generator.tech/)
- Integrate with CI/CD pipelines for API contract testing

For more information about Scalar and its features, visit the [Scalar documentation](https://guides.scalar.com/scalar/introduction).

## Database Schema

### Core Entities
- **User** - Climber profiles with authentication
  - Username, email, password hash
  - Created/updated timestamps

- **Session** - User sessions for authentication
  - Session ID, user ID, expiration
  - Revocation status

- **Crag** - Outdoor climbing areas
  - Name, location, description

- **Wall** - Climbing walls within crags
  - Associated with Crag (many-to-one)

- **Route** - Individual climbing routes
  - Associated with Wall (many-to-one)
  - Grade, type, GPS coordinates, ratings

- **Gym** - Indoor climbing gyms
  - Name, location, facilities (bouldering, top rope, lead, etc.)
  - Contact info, pricing, hours

- **Climb** - Individual climb logs
  - User, Route (outdoor) or Gym (indoor)
  - Climb type (indoor/outdoor)
  - Grade, date, completed status
  - Attempts, falls, rating, notes

### Migrations
Database migrations are handled automatically by GORM on startup. Models are defined in the `models/` directory and registered in `internal/db/postgres.go`.

## Development Workflow

### Adding New Features

1. **Define the feature** - Check Core MVP Features list
2. **Create database models** - Add to `models/` and register in `internal/db/postgres.go`
3. **Create DTOs** - Request/response models in `models/*_dto.go`
4. **Implement handlers** - Business logic in `internal/services/<domain>/`
5. **Register routes** - Wire up in `internal/routes/`
6. **Test locally** - `make run` and test with curl/Postman
7. **Run checks** - `make pre-commit`
8. **Commit** - Follow conventional commit format

### Code Style
- Format code before committing: `make fmt`
- Imports are automatically sorted with local package prefix
- Use structured logging with Zap (not `fmt.Println`)
- Follow Go idioms and best practices

## Infrastructure

### Terraform Modules

Located in `infra/terraform/modules/`:

- **vpc-network** - VPC, subnets, internet gateway, NAT gateway
- **rds-postgresql-db** - RDS PostgreSQL database instance
- **ecr-repository** - Docker image registry
- **alb-ecs** - Application Load Balancer with target groups and security groups
- **ecs-service** - ECS Fargate cluster, service, and task definitions
- **lambda-update-task-ip** - Lambda to associate Elastic IP with ECS tasks
- **elastic-ip** - Static IP allocation
- **ec2-bastion-host** - Bastion host for VPC debugging

### Deployment Commands

```bash
# Format Terraform code
make tf-fmt

# Plan infrastructure changes
make tf-plan

# Apply infrastructure changes
make tf-apply

# Build and push Docker image to ECR
make ecr-deploy

# Deploy to ECS (force new deployment)
make ecs-deploy

# Complete deployment pipeline (ECR + ECS)
make deploy

# Check ECS service status
make ecs-status

# View ECS logs
make ecs-logs
```

### Environment Variables

**Local Development (.env file):**
```
PORT=3000
DB_HOST=localhost
DB_PORT=5432
DB_USER=cruxadmin
DB_PASSWORD=cruxdbpassword
DB_NAME=cruxdb
DB_SSLMODE=disable
JWT_SECRET=your-secret-key-here
```

**Production (ECS Task Definition):**
- Configured via Terraform in `infra/terraform/api.tf`
- Database credentials managed separately (consider AWS Secrets Manager)

## Testing

```bash
# Run all tests with coverage
make test

# Run specific test
go test -v ./internal/services/users/...

# Run with race detection
go test -race ./...
```

## Roadmap

### Completed ✅
- [x] Local development environment
- [x] Core API structure with Fiber
- [x] PostgreSQL database with GORM
- [x] User authentication (JWT + bcrypt)
- [x] Session management with refresh tokens
- [x] Climb logging API (indoor/outdoor)
- [x] Gym model and integration
- [x] AWS infrastructure with Terraform
- [x] ECS Fargate deployment
- [x] Application Load Balancer
- [x] RDS PostgreSQL
- [x] Docker containerization
- [x] Interactive API documentation (OpenAPI + Scalar)

### In Progress 🚧
- [ ] Natural language route search
- [ ] Community features (comments, following, partner matching)
- [ ] Performance analytics and climb statistics
- [ ] Route recommendations

### Planned 📋
- [ ] Mobile app integration
- [ ] Real-time notifications
- [ ] Photo uploads for climbs and routes
- [ ] Weather integration for route conditions
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] GraphQL endpoint (optional)

## Contributing

1. Install pre-commit hooks: `pre-commit install`
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make changes and test: `make test`
4. Format and lint: `make pre-commit`
5. Commit with conventional format: `git commit -m "feat: add route search"`
6. Submit a pull request

### Commit Message Format
```
<type>(<scope>): <subject>

Types: feat, fix, docs, style, refactor, test, chore
Example: feat(search): add natural language route search
```

## Troubleshooting

### Common Issues

**Port 3000 already in use:**
```bash
# Find process using port
lsof -i :3000
# Kill process
kill -9 <PID>
```

**Database connection errors:**
```bash
# Check if db is running
make status
# View db logs
make logs-db
# Reset db
make reset
```

**Cannot connect to Docker:**
```bash
# Ensure Docker is running
docker ps
# Restart Docker Desktop
```

## Project Structure

```
crux-backend/
├── main.go                 # Application entry point
├── models/                 # Database models and DTOs
│   ├── user.go            # User model
│   ├── user_dto.go        # User DTOs
│   ├── auth_dto.go        # Auth DTOs
│   ├── climb.go           # Climb model
│   ├── climb_dto.go       # Climb DTOs
│   ├── gym.go             # Gym model
│   ├── route.go           # Route model
│   ├── wall.go            # Wall model
│   └── crag.go            # Crag model
├── internal/
│   ├── db/                # Database connection & migrations
│   ├── routes/            # HTTP route definitions
│   ├── services/          # Business logic handlers
│   │   ├── auth/         # Authentication
│   │   ├── users/        # User management
│   │   ├── climbs/       # Climb logging
│   │   └── docs/         # API documentation
│   ├── utils/            # Shared utilities
│   │   ├── logger.go     # Zap logger
│   │   ├── jwt.go        # JWT utilities
│   │   └── response.go   # API response helpers
│   └── aws/              # AWS service clients
├── docs/                 # API Documentation
│   ├── openapi.yaml      # OpenAPI 3.0.3 specification
│   └── README.md         # Documentation guide
├── infra/                # Infrastructure as Code
│   ├── terraform/        # Terraform modules
│   │   ├── modules/     # Reusable modules
│   │   ├── api.tf       # API infrastructure
│   │   ├── db.tf        # Database infrastructure
│   │   └── network.tf   # Network infrastructure
│   └── scripts/         # Deployment scripts
├── docker-compose.yml   # Local development stack
├── Dockerfile.dev       # Production Docker image
├── Dockerfile.local     # Development Docker image
├── Makefile            # Development commands
├── go.mod              # Go dependencies
└── .air.toml           # Live reload config
```

## License

[Add your license here]

---

**Status:** 🚧 Active Development
**Current Focus:** Natural language route search MVP

For questions or issues, please open a GitHub issue.
