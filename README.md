# CruxBackend

A high-performance rock climbing API built with Go and [Fiber](https://gofiber.io/) that powers web and mobile applications for the climbing community.

## Overview

CruxBackend serves as the central API for climbers to discover outdoor routes, log gym sessions, track progress, connect with partners, and build a community through route reviews and location sharing.

**Core Features:**
- 🧗 Route discovery and outdoor climbing locations
- 📊 Climb logging and progress tracking
- 👥 Partner matching based on climbing metrics
- ⭐ Route reviews and ratings
- 📍 New location submissions
- 📈 Performance analytics

## Tech Stack

- **API Framework:** [Go Fiber](https://gofiber.io/) - Fast, Express-inspired web framework
- **Database:** [PostgreSQL](https://www.postgresql.org/) - Relational data store
- **Containerization:** [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- **Infrastructure (Planned):** [AWS ECS Fargate](https://aws.amazon.com/fargate/) + [Terraform](https://www.terraform.io/)

## Architecture

```
┌─────────────────┐
│  Web/Mobile App │
└────────┬────────┘
         │
    ┌────▼─────┐
    │   Fiber  │ (API Container)
    │    API   │
    └────┬─────┘
         │
    ┌────▼─────┐
    │ Postgres │ (DB Container)
    └──────────┘
```

Currently runs as a Docker Compose stack locally. Future deployment will use AWS ECS Fargate with Terraform-managed infrastructure.

## Quick Start

### Prerequisites
- [Go 1.24+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.gnu.org/software/make/)

### Local Development

```bash
# Start services (API + Database)
make up

# Run API locally (auto-starts DB)
make run

# View logs
make logs

# Run tests
make test

# Format & lint code
make fmt
make lint

# Database shell
make db-shell

# View all commands
make help
```

The API will be available at `http://localhost:3000`

### Health Check
```bash
curl http://localhost:3000/health
```

## Database Schema

Core entities:
- **Users** - Climber profiles and authentication
- **Routes** - Outdoor and gym climbing routes
- **Climbs** - Individual climb logs and attempts
- **Sessions** - Gym session tracking
- **Walls** - Indoor wall/gym locations

## Development Tools

- **Linting:** [golangci-lint](https://golangci-lint.run/)
- **Pre-commit Hooks:** [pre-commit](https://pre-commit.com/)
- **Live Reload:** [Air](https://github.com/cosmtrek/air)

## Roadmap

- [x] Local development environment
- [x] Core API structure
- [ ] Authentication & authorization
- [ ] Route discovery endpoints
- [ ] Climb logging system
- [ ] Partner matching algorithm
- [ ] Cloud deployment (AWS ECS)
- [ ] Terraform infrastructure
- [ ] CI/CD pipeline

## Contributing

1. Install pre-commit hooks: `pre-commit install`
2. Create a feature branch
3. Make changes and run tests: `make test`
4. Format and lint: `make fmt && make lint`
5. Submit a pull request

## License

[Add your license here]

---

**Status:** 🚧 Active Development
