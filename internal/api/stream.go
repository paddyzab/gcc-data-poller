package api

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/paddyzab/gcc-data-poller/internal/config"
	"github.com/paddyzab/gcc-data-poller/internal/processor"
	"github.com/paddyzab/gcc-data-poller/pkg/pb"
)

type StreamServer struct {
	cfg       *config.Config
	processor *processor.SignalProcessor
	server    *grpc.Server
}

func NewStreamServer(cfg *config.Config, proc *processor.SignalProcessor) *StreamServer {
	return &StreamServer{
		cfg:       cfg,
		processor: proc,
	}
}

func (s *StreamServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.server = grpc.NewServer()
	pb.RegisterSignalStreamerServer(s.server, s)

	log.Printf("gRPC streaming server listening on :%s", s.cfg.GRPCPort)
	return s.server.Serve(lis)
}

func (s *StreamServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}

// SubscribeSignals is the gRPC streaming RPC endpoint
func (s *StreamServer) SubscribeSignals(req *pb.SubscribeRequest, stream pb.SignalStreamer_SubscribeSignalsServer) error {
	log.Printf("Client subscribed to signals (filter: %s)", req.ApiFilter)

	ch := s.processor.Subscribe()
	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			log.Println("Client disconnected")
			return ctx.Err()
		case sig := <-ch:
			// Apply optional filter
			if req.ApiFilter != "" && req.ApiFilter != sig.API {
				continue
			}

			// Convert internal model to pb model
			pbSig := &pb.ActionSignal{
				Api:     sig.API,
				Level:   string(sig.Level),
				Message: sig.Message,
			}

			if err := stream.Send(pbSig); err != nil {
				log.Printf("Error sending signal to client: %v", err)
				return err
			}
		}
	}
}
