# GCC Data Poller PoC

This is the initial project structure for the Proof of Concept (PoC) of a data poller integrating with Google Cloud APIs. 

## Overview
The service performs the following tasks:
1. **Polls Metrics**: Retrieves dummy metrics (latency, error rates) for Firebase Cloud Messaging (FCM) and Identity Toolkit API.
2. **Processes Signals**: Analyzes metrics and generates actionable signals (INFO, WARNING, CRITICAL).
3. **Exposes Downstream API**: Broadcasts these signals to downstream services via a gRPC streaming endpoint.

## Project Structure
- `cmd/server/main.go`: The main entrypoint. Starts the poller and the gRPC server.
- `internal/config`: Loads application configuration via environment variables.
- `internal/models`: Shared types between components.
- `internal/poller`: Handles polling the GCP metrics.
- `internal/processor`: Contains the business logic to generate action signals.
- `internal/api`: Implements the gRPC streaming server.
- `api/proto`: Contains the gRPC protocol buffers definitions.
- `pkg/pb`: Contains the generated gRPC code (currently mocked for PoC).

## Running the PoC

### Prerequisites
- Go 1.21+

### Configuration
The service loads configuration either from environment variables or from a `.config.json` file in the root directory. 
A sample `.config.json` has been provided. **Note:** `.config.json` is ignored by Git to avoid leaking credentials.

### Build and Run
```bash
# Fetch dependencies
make tidy

# Run the server
make run
```
The server will start the poller in the background and expose the gRPC server on port `50051`.

## Generating gRPC Code
A mock implementation is provided in `pkg/pb/signals_mock.go` to allow the PoC to compile without `protoc`. 

If you want to generate the actual gRPC code:
1. Install `protoc` (Protocol Buffers compiler).
2. Install Go plugins:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
3. Run `make proto`
