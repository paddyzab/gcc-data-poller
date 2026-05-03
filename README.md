# GCC Data Poller PoC

This is the initial project structure for the Proof of Concept (PoC) of a data poller integrating with Google Cloud APIs. 

## Overview
The service performs the following tasks:
1. **Polls Metrics**: Retrieves dummy metrics (latency, error rates) for Firebase Cloud Messaging (FCM) and Identity Toolkit API.
2. **Processes Signals**: Analyzes metrics and generates actionable signals (INFO, WARNING, CRITICAL).
3. **Exposes Downstream API**: Broadcasts these signals to downstream services via a gRPC streaming endpoint.

## Project Structure
- `cmd/server/main.go`: The main entrypoint. Starts the poller and the gRPC server.
- `cmd/client/main.go`: A demo client to connect to the gRPC server and receive signals.
- `internal/config`: Loads application configuration via environment variables.
- `internal/models`: Shared types between components.
- `internal/poller`: Handles polling the GCP metrics.
- `internal/processor`: Contains the business logic to generate action signals.
- `internal/api`: Implements the gRPC streaming server.
- `api/proto`: Contains the gRPC protocol buffers definitions.
- `pkg/pb`: Contains the generated gRPC code from `protoc`.

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

### Integrating with the Server (PoC)

To verify the integration or connect a downstream service to this PoC, follow these steps:

#### 1. Run the Demo Client
A demo client is included in the project to instantly test the streaming connection. In a separate terminal, run:
```bash
go run ./cmd/client
```
The client will connect to `localhost:50051`, subscribe to "FCM" signals, and stream them to your terminal in real-time.

#### 2. Connect Your Own Services
To integrate a different downstream Go service with this poller:
1. Copy the `api/proto/signals.proto` file into your downstream service's repository.
2. Compile the protobuf file in your repository to generate the Go client stubs using `protoc`.
3. Dial the gRPC connection and use `pb.NewSignalStreamerClient(conn)` to initiate a subscription stream, just as demonstrated in the `cmd/client/main.go` file.

## Generating gRPC Code
If you modify `api/proto/signals.proto`, you will need to regenerate the gRPC code:

1. Install `protoc` (Protocol Buffers compiler).
2. Install Go plugins:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
```
3. Run `make proto`
