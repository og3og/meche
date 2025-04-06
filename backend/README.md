# Backend gRPC and HTTP API Service

This service provides a gRPC API with an HTTP gateway, allowing the same service to be accessed via both gRPC and REST endpoints.

## Prerequisites

- Go 1.24.2 or later
- [buf](https://buf.build/) for protocol buffer management
- `curl` for testing HTTP endpoints (optional)
- `jq` for formatting JSON responses (optional)

## Quick Start

1. **Install all dependencies and run the server:**
   ```bash
   make run-all
   ```
   This will:
   - Clean any existing builds
   - Install required dependencies
   - Generate protocol buffer code
   - Build the server and client
   - Start the server

2. **Test the endpoints:**

   - Using HTTP (REST):
     ```bash
     make test-http
     ```
     Or with curl directly:
     ```bash
     curl -X POST -H "Content-Type: application/json" \
       -d '{"name": "World"}' \
       http://localhost:8080/v1/greeter/hello
     ```

   - Using gRPC client:
     ```bash
     make client
     ```

## Available Endpoints

- gRPC: `localhost:50051`
- HTTP: `http://localhost:8080/v1/greeter/hello` (POST)

## Available Make Commands

- `make run-all` - Clean, build, and run everything
- `make proto` - Generate protobuf code
- `make build` - Build the server and client
- `make run` - Run the server
- `make client` - Run the client
- `make clean` - Clean build artifacts
- `make deps` - Install dependencies
- `make test-http` - Test the HTTP endpoint
- `make help` - Show all available commands

## Service Definition

The service provides a simple greeting endpoint:

```protobuf
service GreeterService {
  rpc SayHello (HelloRequest) returns (HelloResponse) {
    option (google.api.http) = {
      post: "/v1/greeter/hello"
      body: "*"
    };
  }
}
```

## Development

To modify the service:

1. Update the proto definitions in `proto/v1/`
2. Generate new code: `make proto`
3. Implement the new service methods in `internal/service/`
4. Rebuild and run: `make run-all` 