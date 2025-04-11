.PHONY: dev build clean generate test

# Environment variables
export GOOGLE_CLIENT_ID=1087371858942-dmij25b6onipcj9irjaq91epju2306lm.apps.googleusercontent.com
export GOOGLE_CLIENT_SECRET=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

# Default target
all: generate build

# Run the application in development mode with hot reload
dev:
	air

# Generate templ files
generate:
	templ generate

# Build the application
build: generate
	go build -o ./tmp/server ./cmd/server

# Run the application
run: build
	./tmp/server

# Clean build artifacts
clean:
	rm -rf ./tmp
	go clean

# Run tests
test:
	go test ./... -v

# Install development dependencies
install-deps:
	go install github.com/cosmtrek/air@latest
	go install github.com/a-h/templ/cmd/templ@latest

# Format code
fmt:
	go fmt ./...
	templ fmt ./...

# Lint code
lint:
	go vet ./...

# Help command to show available commands
help:
	@echo "Available commands:"
	@echo "  make dev          - Run the application in development mode with hot reload"
	@echo "  make generate     - Generate templ files"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the built application"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make fmt          - Format Go and Templ files"
	@echo "  make lint         - Run Go linter"
	@echo "  make install-deps - Install development dependencies" 