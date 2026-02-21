# ---------------------------------------
# Project Makefile
# ---------------------------------------
# Usage: `make help` to list targets.
# Put a short description after `##` on each target to show up in help.

SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c

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
        start stop restart logs

## help: Display this help message with a list of all targets and their usage.
help:
	@echo "Usage: make <target>"
	@echo
	@echo "Available targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/  /'

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

## start: Start development services (e.g., MySQL) in detached mode
start:
	@echo "→ Starting development services..."
	@$(DOCKER_COMPOSE) up -d

## stop: Stop development services
stop:
	@echo "→ Stopping development services..."
	@$(DOCKER_COMPOSE) down

## restart: Restart development services
restart:
	@echo "→ Restarting development services..."
	@$(DOCKER_COMPOSE) down
	@$(DOCKER_COMPOSE) up -d

## logs: Tail service logs
logs:
	@$(DOCKER_COMPOSE) logs -f

## clean: Remove build artifacts and coverage output
clean:
	@echo "→ Cleaning..."
	@rm -rf bin "$(COVERAGE_OUTPUT)"
