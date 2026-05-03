.PHONY: build run test tidy proto

build:
	go build -o bin/poller ./cmd/server

run: build
	./bin/poller

test:
	go test ./...

tidy:
	go mod tidy

proto:
	@echo "Requires protoc and grpc-go plugins"
	protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           api/proto/signals.proto
