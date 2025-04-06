# Backend gRPC and HTTP API Service

This service provides a gRPC API with an HTTP gateway, allowing the same service to be accessed via both gRPC and REST endpoints. The service is secured using Auth0 authentication.

## Prerequisites

- Go 1.24.2 or later
- [buf](https://buf.build/) for protocol buffer management
- `curl` for testing HTTP endpoints (optional)
- `jq` for formatting JSON responses (optional)
- An Auth0 account and configured API

## Auth0 Setup

1. Create a new API in your Auth0 dashboard:
   - Name: Choose a name for your API
   - Identifier: This will be your `AUTH0_AUDIENCE`
   - Signing Algorithm: RS256

2. Copy your Auth0 domain (e.g., `your-domain.auth0.com`) and API identifier.

3. Create a `.env` file in the backend directory:
   ```bash
   cp .env.example .env
   ```

4. Update the `.env` file with your Auth0 configuration:
   ```
   AUTH0_DOMAIN=your-domain.auth0.com
   AUTH0_AUDIENCE=your-api-identifier
   ```

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
     # Get a token from Auth0
     TOKEN="your-auth0-token"
     
     # Make an authenticated request
     curl -X POST \
       -H "Authorization: Bearer $TOKEN" \
       -H "Content-Type: application/json" \
       -d '{"name": "World"}' \
       http://localhost:8080/v1/greeter/hello
     ```

   - Using gRPC client:
     ```bash
     # The client needs to be updated to include the token
     make client
     ```

## Available Endpoints

- gRPC: `localhost:50051`
- HTTP: `http://localhost:8080/v1/greeter/hello` (POST)

All endpoints require a valid Auth0 JWT token in the Authorization header.

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

## Authentication

The service uses Auth0 for authentication. Each request must include a valid JWT token in the Authorization header:

```
Authorization: Bearer your-auth0-token
```

For development purposes, you can disable authentication by setting:
```
AUTH_DISABLED=true
```
in your `.env` file. 