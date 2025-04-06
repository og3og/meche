.PHONY: all build run clean proto client run-all test-http

# Default target
all: proto build run

# Run everything (clean, build, and run)
run-all:
	@echo "Cleaning up..."
	@pkill -f "bin/server" || true
	@lsof -ti:8080 | xargs kill -9 2>/dev/null || true
	@make clean
	@echo "Installing dependencies..."
	@make deps
	@echo "Generating proto files..."
	@make proto
	@echo "Building..."
	@make build
	@echo "Starting server..."
	@make run

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	buf generate

# Build the server and client
build:
	@echo "Building server..."
	go build -o bin/server cmd/server/main.go
	@echo "Building client..."
	go build -o bin/client cmd/client/main.go

# Run the server
run:
	@echo "Starting gRPC server..."
	./bin/server

# Run the client
client:
	@echo "Running gRPC client..."
	./bin/client

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f proto/v1/*.pb.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.0

# Test the HTTP endpoint
test-http:
	@echo "Testing HTTP endpoint..."
	@if [ -z "$$AUTH0_TOKEN" ] && [ "$$(grep -E '^AUTH_DISABLED=true$$' .env)" = "" ]; then \
		echo "Error: Please set AUTH0_TOKEN environment variable or enable AUTH_DISABLED in .env"; \
		exit 1; \
	fi
	@if [ -n "$$AUTH0_TOKEN" ]; then \
		curl -X POST -v \
			-H "Authorization: Bearer $$AUTH0_TOKEN" \
			-H "Content-Type: application/json" \
			-d '{"name": "World"}' \
			http://localhost:8080/v1/greeter/hello ; \
	else \
		curl -X POST \
			-H "Content-Type: application/json" \
			-d '{"name": "World"}' \
			http://localhost:8080/v1/greeter/hello | jq .; \
	fi

# Help command
help:
	@echo "Available commands:"
	@echo "  make run-all  - Clean, build, and run everything"
	@echo "  make proto    - Generate protobuf code"
	@echo "  make build    - Build the server and client"
	@echo "  make run      - Run the server"
	@echo "  make client   - Run the client"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make deps     - Install dependencies"
	@echo "  make test-http- Test the HTTP endpoint (requires AUTH0_TOKEN or AUTH_DISABLED=true)"
	@echo "  make help     - Show this help message"
