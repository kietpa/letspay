# Locations
DEPLOY := deploy/docker-compose
ENV_FILE := .env
MIGRATIONS_DIR := migrations

# Database URL (can be set via .env or directly passed to goose)
DATABASE_URL ?= postgres://letsuser:letspassword@localhost:5372/letspay?sslmode=disable

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build-all        - Build all services"
	@echo "  make run              - Run all services with Docker Compose"
	@echo "  make stop             - Stop all running containers"
	@echo "  make tidy             - Run go mod tidy in all services and pkg"
	@echo "  make migrate-up       - Run all DB migrations"
	@echo "  make migrate-down     - Rollback the last DB migration"
	@echo "  make status           - Show current migration status"

# -----------------------------------------------------------------------------
# BUILD
# -----------------------------------------------------------------------------
.PHONY: build-all
build-all:
	@echo "🔨 Building all services..."
	cd services/user && go build -o user-service
	cd services/payment && go build -o payment-service
	cd services/api-gateway && go build -o api-gateway

# -----------------------------------------------------------------------------
# DOCKER COMPOSE
# -----------------------------------------------------------------------------
.PHONY: run
run:
	@echo "🚀 Running services via Docker Compose..."
	docker-compose -f $(DEPLOY)/docker-compose.yml --env-file $(ENV_FILE) up --build

.PHONY: stop
stop:
	@echo "🛑 Stopping containers..."
	docker-compose -f $(DEPLOY)/docker-compose.yml down

# -----------------------------------------------------------------------------
# DB MIGRATIONS (using goose CLI)
# -----------------------------------------------------------------------------
.PHONY: migrate-up
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

.PHONY: status
status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status

# -----------------------------------------------------------------------------
# UTILITIES
# -----------------------------------------------------------------------------
.PHONY: tidy
tidy:
	@echo "🧹 Running go mod tidy..."
	cd services/user && go mod tidy
	cd services/payment && go mod tidy
	cd services/api-gateway && go mod tidy
	cd pkg && go mod tidy
