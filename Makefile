# ---------------------------------------
# Project Makefile
# ---------------------------------------
# Usage: `make help` to list targets.
# Put a short description after `##` on each target to show up in help.


# The name of the binary or project directory
PROJECT_NAME = snippetbox

# --------------------------------------------------------------------
# Shell Configuration
# --------------------------------------------------------------------

# Use bash instead of the default /bin/sh for better feature support
SHELL        := bash

# .ONESHELL ensures all lines in a recipe run in a single shell instance.
# This allows variables and directory changes to persist between lines.
.ONESHELL:

# -e: Exit immediately if a command fails.
# -u: Treat unset variables as an error.
# -o pipefail: Ensure the exit code of a pipeline is the status of the last command to exit with a non-zero status.
.SHELLFLAGS  := -eu -o pipefail -c

# --------------------------------------------------------------------
# Database migrations
# --------------------------------------------------------------------

DATABASE_DSN     ?= dev:demo@tcp(127.0.0.1:3306)/snippetbox
MIGRATE          := migrate
MIGRATION_FOLDER := ./migrations

# --------------------------------------------------------------------
# Tools / Config
# --------------------------------------------------------------------

GO              ?= go
GOFLAGS         ?=
RACE            ?= -race
PORT            ?= 4000

# Docker
COMPOSE_FILE    ?= compose.yml
DOCKER_COMPOSE  ?= docker compose -f $(COMPOSE_FILE)

# Test config
COVERAGE_OUTPUT ?= cover.out
GOTEST_FLAGS    ?= -coverprofile=$(COVERAGE_OUTPUT)

.PHONY: build clean deps fmt vet local \
        test test-race test-integration help \
        start stop restart logs \
        migration-create migration-up migration-down migration-status migration-force

# --------------------------------------------------------------------
# Help
# --------------------------------------------------------------------

## help: Display this help message with a list of all targets and their usage.
help:
	@echo "Usage: make <target>"
	@echo
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/  /'

# --------------------------------------------------------------------
# Go project tasks
# --------------------------------------------------------------------

## deps: Manage Go dependencies (tidy + download)
deps:
	@echo "→ Tidying and downloading modules..."
	@$(GO) mod tidy
	@$(GO) mod download

## fmt: Format source code
fmt:
	@echo "→ Formatting code..."
	@$(GO) fmt ./...

## vet: Run go vet (static analysis)
vet:
	@echo "→ Vetting..."
	@$(GO) vet ./...

## build: Build the Go server binary
build: deps fmt vet
	@echo "→ Building application..."
	@mkdir -p bin
	@$(GO) build $(GOFLAGS) -o bin/snippetbox ./cmd/web

## local: Start the Go server locally
local: build
	@echo "→ Starting Go server on :$(PORT)"
	@./bin/snippetbox -addr=":$(PORT)"

## test: Run unit tests with coverage
test: deps fmt vet
	@echo "→ Running unit tests..."
	@$(GO) test $(GOFLAGS) $(GOTEST_FLAGS) -short ./...

## test-race: Run tests with the race detector
test-race: deps fmt vet
	@echo "→ Running tests with race detector..."
	@$(GO) test $(GOFLAGS) $(RACE) $(GOTEST_FLAGS) ./...

## test-integration: Run integration tests (serial)
test-integration: deps fmt vet
	@echo "→ Running integration tests..."
	@$(GO) test $(GOFLAGS) -p 1 $(GOTEST_FLAGS) ./...

## clean: Remove build artifacts and coverage output
clean:
	@echo "→ Cleaning..."
	@rm -rf bin "$(COVERAGE_OUTPUT)"

# --------------------------------------------------------------------
# Docker (services)
# --------------------------------------------------------------------

## start: Start development services (e.g., MySQL) in detached mode
start:
	@echo "→ Starting development services..."
	@$(DOCKER_COMPOSE) up -d

## stop: Stop development services and remove volumes (fresh start next time)
stop:
	@echo "→ Stopping development services..."
	@$(DOCKER_COMPOSE) down

## restart: Restart development services without removing volumes
restart: stop start

## prune: Deep clean of services, volumes, and orphaned images
prune:
	@echo "→ Deep cleaning development environment..."
	@$(DOCKER_COMPOSE) down -v --remove-orphans
	@docker image prune -f --filter "label=com.docker.compose.project=$(PROJECT_NAME)"

## logs: Tail service logs
logs:
	@$(DOCKER_COMPOSE) logs -f

# --------------------------------------------------------------------
# Migrations
# --------------------------------------------------------------------

## migration-create: Create a new SQL migration. Pass NAME=<migration_name>
migration-create:
	@if [ -z "$(NAME)" ]; then \
		echo "ERROR: NAME is required. Usage: make migration-create NAME=your_migration_name"; \
		exit 1; \
	fi
	@echo "→ Creating new migration: $(NAME)"
	@$(MIGRATE) create -ext sql -dir $(MIGRATION_FOLDER) -seq $(NAME)
	@echo "→ Migration $(NAME) created successfully."

## migration-up: Apply all pending migrations
migration-up:
	@echo "→ Applying all migrations from $(MIGRATION_FOLDER)..."
	@$(MIGRATE) -path $(MIGRATION_FOLDER) -database "mysql://$(DATABASE_DSN)" up
	@echo "→ All migrations applied successfully."

## migration-down: Rollback migrations interactively
migration-down:
	@echo "→ Rolling back migrations..."
	@read -p "Number of migrations to rollback (default: 1): " NUM; \
	NUM=$${NUM:-1}; \
	$(MIGRATE) -path $(MIGRATION_FOLDER) -database "mysql://$(DATABASE_DSN)" down $${NUM}; \
	echo "→ Rolled back $${NUM} migration(s)."

## migration-status: Show migration status
migration-status:
	@echo "→ Showing migration status..."
	@$(MIGRATE) -path $(MIGRATION_FOLDER) -database "mysql://$(DATABASE_DSN)" version

## migration-force: Force the database to a specific migration version
migration-force:
	@echo "→ Forcing migration version..."
	@read -p "Enter the version to force: " VERSION; \
	if [ -z "$${VERSION}" ]; then \
		echo "ERROR: Version is required"; exit 1; \
	fi; \
	$(MIGRATE) -path $(MIGRATION_FOLDER) -database "mysql://$(DATABASE_DSN)" force $${VERSION}; \
	echo "→ Migration version forced to $${VERSION}."
