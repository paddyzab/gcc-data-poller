// Code generated for PoC purposes. In a real environment, generate this using protoc.
package pb

import (
	"context"
	"google.golang.org/grpc"
)

type SubscribeRequest struct {
	ApiFilter string
}

type ActionSignal struct {
	Api       string
	Level     string
	Message   string
}

type SignalStreamerServer interface {
	SubscribeSignals(*SubscribeRequest, SignalStreamer_SubscribeSignalsServer) error
}

type SignalStreamer_SubscribeSignalsServer interface {
	Send(*ActionSignal) error
	Context() context.Context
	grpc.ServerStream
}

func RegisterSignalStreamerServer(s *grpc.Server, srv SignalStreamerServer) {
	// Mock implementation. Real protoc generates exact server registration.
}
