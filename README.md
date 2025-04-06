# Meche Monorepo

This monorepo contains both the backend gRPC/HTTP API service and the frontend web application.

## Project Structure

```
.
├── backend/           # Backend service
│   ├── cmd/
│   │   ├── server/   # gRPC and HTTP server implementation
│   │   └── client/   # gRPC client implementation
│   ├── internal/
│   │   └── service/  # Service implementation
│   ├── proto/
│   │   └── v1/      # Protocol buffer definitions
│   ├── bin/         # Compiled binaries
│   ├── buf.yaml     # Buf configuration
│   ├── buf.gen.yaml # Buf generation configuration
│   ├── go.mod       # Go module file
│   └── Makefile     # Backend build commands
└── www/             # Frontend web application (coming soon)
```

## Components

### Backend Service

The backend provides a gRPC service with an HTTP gateway, allowing the same service to be accessed via both gRPC and REST endpoints. [More details](backend/README.md)

### Frontend Application

The frontend web application will be a modern web interface that communicates with the backend service. (Coming soon)

## Development

### Backend Development

Navigate to the backend directory and follow the instructions in its README:
```bash
cd backend
make run-all  # Builds and runs the backend service
```

### Frontend Development

Frontend development instructions will be added soon. 