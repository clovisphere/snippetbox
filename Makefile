# ---------------------------------------
# Project Makefile
# ---------------------------------------
# Usage: `make help` to list targets.
# Put a short description after `##` on each target to show up in help.

SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c

# Tools / Config
GO              ?= go
GOFLAGS         ?=
GOTEST_FLAGS    ?= -coverprofile=$(COVERAGE_OUTPUT)
RACE            ?= -race
COVERAGE_OUTPUT ?= cover.out

# Port is the port the application uses (default to 4000)
PORT ?= 4000

.PHONY: build clean deps fmt vet local \
	    report test test-race test-integration ci help generate

## help: Show this help
help:
	@printf "%b\n" "$(BLUE)Usage: make <target>$(NC)"
	@awk '/^## / {                                   \
		line = substr($$0, 4);                       \
		i = index(line, ":");                        \
		if (i) {                                     \
			name = substr(line, 1, i-1);             \
			desc = substr(line, i+1);                \
			gsub(/^[ \t]+/, "", desc);               \
			printf "  \033[36m%-20s\033[0m %s\n",    \
			       name, desc;                       \
		}                                            \
	}' $(MAKEFILE_LIST)

## deps: Manage Go dependencies (tidy + download)
deps:
	@echo "→ Tidying and downloading modules..."
	@$(GO) mod tidy
	@$(GO) mod download

## fmt: Format source code
fmt:
	@echo "→ Formatting code..."
	@$(GO) fmt ./...
	@echo "→ gofumpt (optional): install with 'go install mvdan.cc/gofumpt@latest' and run manually if desired."

## vet: Run go vet (static analysis)
vet:
	@echo "→ Vetting..."
	@$(GO) vet ./...

## build: Build the Go server binary
build: deps fmt vet
	@echo "→ Building application..."
	@$(GO) build $(GOFLAGS) -o bin/snippetbox ./cmd/web

## local: Start the Go server locally
local: build
	@echo "→ Starting Go server..."
	@echo ""
	@./bin/snippetbox -addr=":$(PORT)"

## generate: Run code generation (if any go:generate directives exist)
generate:
	@echo "→ Running go generate (if present)..."
	@$(GO) generate ./...

## test: Run unit tests with coverage
test: deps fmt vet
	@echo "→ Running unit tests..."
	@$(GO) test $(GOFLAGS) $(GOTEST_FLAGS) -short ./...

## test-race: Run tests with the race detector
test-race: deps fmt vet
	@echo "→ Running tests (race)..."
	@$(GO) test $(GOFLAGS) $(RACE) $(GOTEST_FLAGS) ./...

## test-integration: Run integration tests (serial)
test-integration: deps fmt vet
	@echo "→ Running integration tests..."
	@$(GO) test $(GOFLAGS) -p 1 $(GOTEST_FLAGS) ./...

## report: Generate HTML coverage report to bin/coverage.html
report:
	@echo "→ Generating coverage report..."
	@mkdir -p bin
	@{ test -f "$(COVERAGE_OUTPUT)" && \
	   $(GO) tool cover -html=$(COVERAGE_OUTPUT) -o bin/coverage.html && \
	   echo "Coverage report at bin/coverage.html"; } || \
	   echo "No coverage file '$(COVERAGE_OUTPUT)' found. Run 'make test' first."

## lint: Run golangci-lint (configure via .golangci.yml or LINT_ARGS)
# Linting notes:
#   - Prefer a pinned .golangci.yml to ensure consistent rules across dev/CI.
#   - Use LINT_ARGS to try rules ad hoc, e.g.:
#       make lint LINT_ARGS='run --enable=gofmt --timeout=5m ./...'
lint:
	@echo "→ Running golangci-lint..."
	@$(GOLANGCI_LINT) $(LINT_ARGS)

## ci: Run the CI pipeline (lint + tests + report)
ci: lint test report
	@echo "→ CI pipeline completed."

## clean: Remove build artifacts and coverage output
clean:
	@echo "→ Cleaning..."
	@rm -rf bin/ "$(COVERAGE_OUTPUT)"
