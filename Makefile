# Makefile for Ethereum Fetcher Project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=ethereum-fetcher

# Directories
BUILD_DIR=build

# Environment variables
ENV_FILE=.env

include .env
export

# PostgreSQL Docker configuration
# Extract connection details from DB_CONNECTION_URL
POSTGRES_USER=$(shell echo $(DB_CONNECTION_URL) | sed -E 's|postgresql://([^:]+):.*|\1|')
POSTGRES_PASSWORD=$(shell echo $(DB_CONNECTION_URL) | sed -E 's|postgresql://[^:]+:([^@]+).*|\1|')
POSTGRES_HOST=$(shell echo $(DB_CONNECTION_URL) | sed -E 's|postgresql://[^@]+@([^:]+).*|\1|')
POSTGRES_PORT=$(shell echo $(DB_CONNECTION_URL) | sed -E 's|.*:([0-9]+)/.*|\1|' || echo 5432)
POSTGRES_DB=$(shell echo $(DB_CONNECTION_URL) | sed -E 's|.*/([^?]+).*|\1|')
POSTGRES_CONTAINER_NAME=ethereum-fetcher-postgres

# Docker configuration
DOCKER_IMAGE_NAME=limeapi
DOCKER_PORT=8080

# Default target
.PHONY: all
all: postgres build-all run

# Build all
.PHONY: build-all
build-all: setup docker-build

# Build the application
.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) main.go

# Run the application
.PHONY: run
run:
	$(GORUN) main.go

# Run tests
.PHONY: test
test:
	$(GOTEST) ./tests/...

# Clean build files
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	$(GOCMD) clean

# Install dependencies
.PHONY: deps
deps:
	$(GOMOD) tidy
	$(GOMOD) download

# Generate documentation (if needed)
.PHONY: docs
docs:
	# Add documentation generation command if applicable

# Lint the code
.PHONY: lint
lint:
	golangci-lint run

# Format the code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Full setup and build
.PHONY: setup
setup: clean deps build

## DATABASE

# Start PostgreSQL in Docker
.PHONY: postgres
postgres:
	@if [ "$$(docker ps -a | grep $(POSTGRES_CONTAINER_NAME))" ]; then \
		docker start $(POSTGRES_CONTAINER_NAME); \
	else \
		docker run --name $(POSTGRES_CONTAINER_NAME) \
			-e POSTGRES_DB=$(POSTGRES_DB) \
			-e POSTGRES_USER=$(POSTGRES_USER) \
			-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
			-p $(POSTGRES_PORT):$(POSTGRES_PORT) \
			-d postgres:15-alpine; \
	fi

# Stop PostgreSQL Docker container
.PHONY: postgres-stop
postgres-stop:
	@docker stop $(POSTGRES_CONTAINER_NAME)

# Restart PostgreSQL Docker container
.PHONY: postgres-restart
postgres-restart: postgres-stop postgres

# Connect to PostgreSQL database
.PHONY: postgres-cli
postgres-cli:
	@docker exec -it $(POSTGRES_CONTAINER_NAME) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

# Docker build and run targets
docker-build:
	docker build -t $(DOCKER_IMAGE_NAME) .

docker-run: docker-build
	docker run -p $(DOCKER_PORT):$(DOCKER_PORT) $(DOCKER_IMAGE_NAME)

docker-stop:
	docker stop $$(docker ps -q --filter ancestor=$(DOCKER_IMAGE_NAME)) 2>/dev/null || true

docker-clean: docker-stop
	docker rmi $(DOCKER_IMAGE_NAME) 2>/dev/null || true

docker-rebuild: docker-clean docker-build

.PHONY: docker-build docker-run docker-stop docker-clean docker-rebuild

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all      - Build and run the application (default)"
	@echo "  build    - Compile the application"
	@echo "  build-all - Compile the application and create docker image"
	@echo "  run      - Run the application"
	@echo "  test     - Run tests"
	@echo "  clean    - Remove build artifacts"
	@echo "  deps     - Download and tidy dependencies"
	@echo "  docs     - Generate documentation"
	@echo "  lint     - Run linter"
	@echo "  fmt      - Format code"
	@echo "  setup    - Install dependencies and build"
	@echo "  postgres - Start PostgreSQL in Docker"
	@echo "  postgres-stop - Stop PostgreSQL Docker container"
	@echo "  postgres-rm - Remove PostgreSQL Docker container"
	@echo "  postgres-restart - Restart PostgreSQL Docker container"
	@echo "  postgres-cli - Connect to PostgreSQL database"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run - Run Docker container"
	@echo "  docker-stop - Stop Docker container"
	@echo "  docker-clean - Remove Docker image"
	@echo "  docker-rebuild - Rebuild Docker image"
	@echo "  help     - Show this help message"