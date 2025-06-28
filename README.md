# WebAssembly (Wasm) Authentication Filter Sample for Envoy

A sample implementation of a Wasm authentication filter for Envoy Proxy written in Go.

## Overview

This project demonstrates how to use Envoy Proxy's Wasm filter capabilities to implement HTTP request authentication.

### Architecture

```
[Client] → [Envoy Proxy + Wasm Filter] → [Backend Service]
```

- **Envoy Proxy**: Acts as a reverse proxy and executes the Wasm filter
- **Wasm Filter**: Authentication filter implemented in Go
- **Backend Service**: Simple HTTP server

### Features

- Bearer token authentication
- Skip authentication for health check endpoints
- Add user information headers upon successful authentication
- Add filter identification headers to responses

## Prerequisites

- Docker
- Docker Compose
- Go 1.24 or later

## Setup Instructions

1. Clone the repository
```bash
git clone https://github.com/keitaj/envoy-wasm-sample.git
cd envoy-wasm-sample
```

2. Build the Wasm filter and start services
```bash
make run
```

This will:
- Build the Go Wasm filter
- Build Docker images
- Start Envoy and Backend services

## Testing

### Individual Test Execution

1. **Health check (no authentication required)**
```bash
curl -i http://localhost:10000/health
```

2. **Access without authentication (401 error)**
```bash
curl -i http://localhost:10000/api/data
```

3. **Access with invalid token (401 error)**
```bash
curl -i -H "Authorization: Bearer invalid-token" http://localhost:10000/api/data
```

4. **Access with valid user token (success)**
```bash
curl -i -H "Authorization: Bearer secret-token-123" http://localhost:10000/api/data
```

5. **Access with valid admin token (success)**
```bash
curl -i -H "Authorization: Bearer admin-token-456" http://localhost:10000/api/data
```

### Run All Tests

```bash
make test
```

## System Behavior

### Authentication Flow

1. **Request Reception**
   - HTTP request from client arrives at Envoy
   - Wasm filter's `OnHttpRequestHeaders` is called

2. **Path Check**
   - `/health` endpoint skips authentication
   - Other paths undergo authentication

3. **Authentication Process**
   - Check for Authorization header
   - Validate Bearer token
   - Valid tokens:
     - `secret-token-123`: user role
     - `admin-token-456`: admin role

4. **On Successful Authentication**
   - Add user information to `x-auth-user` header
   - Forward request to backend service

5. **On Authentication Failure**
   - Return 401 Unauthorized response
   - Return error message in JSON format

### Response Processing

All responses include the `x-wasm-filter: go-auth` header, confirming that the Wasm filter is functioning.

### Logging

The Wasm filter outputs the following logs:
- On plugin start: "plugin started"
- On request processing: path information
- On authentication success/failure: detailed information

These logs can be viewed in Envoy's logs:
```bash
docker logs envoy
```

## Directory Structure

```
.
├── Makefile              # Build, run, and test command definitions
├── docker-compose.yaml   # Docker Compose configuration
├── envoy.yaml           # Envoy configuration file
├── filter.wasm          # Built Wasm file (auto-generated)
├── backend/             # Backend service
│   ├── Dockerfile
│   └── main.go
└── wasm-filter/         # Wasm filter source code
    ├── go.mod
    ├── go.sum
    └── main.go
```

## Troubleshooting

### If Wasm Filter Fails to Load

1. Check Envoy logs
```bash
docker logs envoy
```

2. Verify Go version (1.24 or later required)
```bash
go version
```

3. Confirm Wasm file is built correctly
```bash
ls -la filter.wasm
```

### Cleanup

Stop services and remove build files:
```bash
make clean
```

## Technical Details

- **Envoy**: v1.34-latest
- **Go**: 1.24 (wasip1/wasm target)
- **proxy-wasm-go-sdk**: Wasm plugin development SDK for Envoy

The Wasm filter sets up the VM context in the `init()` function and creates an HTTP context for each request to handle processing.