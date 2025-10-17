# Development Dockerfile with hot-reload support
# ==============================================
# This Dockerfile is optimized for local development with Air hot-reload.
# For production, use Dockerfile.prod instead.

FROM golang:1.24-alpine

# Install build dependencies and tools
RUN apk add --no-cache \
    git \
    wget \
    curl \
    make

# Install Air for hot-reload (compatible with Go 1.24)
# Using v1.49.0 - last stable version before repository migration
RUN go install github.com/cosmtrek/air@v1.49.0

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for better caching
# Docker will cache this layer unless dependencies change
COPY go.mod go.sum ./

# Download dependencies (cached unless go.mod/go.sum changes)
RUN go mod download && go mod verify

# Copy Air configuration
COPY .air.toml ./

# Copy source code
# Note: In docker-compose.yml, we mount the source as a volume
# so changes are reflected immediately without rebuild
COPY . .

# Create tmp directory for Air builds
RUN mkdir -p /app/tmp

# Expose the application port
ARG PORT=3000
EXPOSE ${PORT}

# Set Air as the entrypoint for hot-reload
ENTRYPOINT ["air", "-c", ".air.toml"]
