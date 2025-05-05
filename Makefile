.PHONY: build run test

PORT ?= 4000

all: run

build:
	@echo "Building binary..."
	@go build -o bin/snippetbox ./cmd/web

run: build
	@echo "Running web server... (development mode)"
	@go run ./cmd/web -addr=":$(PORT)"

test:
	@echo "Running tests..."
	@go test -v ./...

clear:
	@echo "Clearing build artifacts..."
	@rm -rf bin/*
