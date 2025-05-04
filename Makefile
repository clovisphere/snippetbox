.PHONY: build run test

all: run

build:
	@echo "Building binary..."
	@go build -o bin/snippetbox ./cmd/web

run: build
	@echo "Running web server... (development mode)"
	@go run ./cmd/web

test:
	@echo "Running tests..."
	@go test -v ./...

clear:
	@echo "Clearing build artifacts..."
	@rm -rf bin/*
