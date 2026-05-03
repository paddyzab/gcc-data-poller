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
	mkdir -p pkg/pb
	protoc --go_out=pkg/pb --go_opt=paths=source_relative \
           --go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
           -I=api/proto api/proto/signals.proto
