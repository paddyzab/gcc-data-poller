package api

import (
	"context"
	"testing"
	"time"

	"github.com/paddyzab/gcc-data-poller/internal/config"
	"github.com/paddyzab/gcc-data-poller/internal/models"
	"github.com/paddyzab/gcc-data-poller/internal/processor"
	"github.com/paddyzab/gcc-data-poller/pkg/pb"
	"google.golang.org/grpc"
)

// mockServerStream implements pb.SignalStreamer_SubscribeSignalsServer
type mockServerStream struct {
	grpc.ServerStream
	ctx      context.Context
	messages []*pb.ActionSignal
}

func (m *mockServerStream) Context() context.Context {
	return m.ctx
}

func (m *mockServerStream) Send(sig *pb.ActionSignal) error {
	m.messages = append(m.messages, sig)
	return nil
}

func TestStreamServer_SubscribeSignals_Integration(t *testing.T) {
	// 1. Setup the environment
	cfg := &config.Config{GRPCPort: "50051"}
	proc := processor.NewSignalProcessor()
	server := NewStreamServer(cfg, proc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := &mockServerStream{
		ctx:      ctx,
		messages: make([]*pb.ActionSignal, 0),
	}

	// Requesting only FCM signals
	req := &pb.SubscribeRequest{
		ApiFilter: "FCM",
	}

	// 2. Start SubscribeSignals in a goroutine (it blocks until context is canceled)
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.SubscribeSignals(req, stream)
	}()

	// Wait a brief moment to ensure the subscription is registered
	time.Sleep(50 * time.Millisecond)

	// 3. Inject raw metrics via the processor
	
	// Metric 1: FCM (Matches filter, Critical latency/error)
	proc.ProcessMetric(models.RawMetric{
		API:       "FCM",
		Timestamp: time.Now(),
		LatencyMs: 600,
		ErrorRate: 0.1,
	})

	// Metric 2: IdentityToolkit (Does NOT match filter, Warning latency)
	proc.ProcessMetric(models.RawMetric{
		API:       "IdentityToolkit",
		Timestamp: time.Now(),
		LatencyMs: 300,
		ErrorRate: 0.02,
	})

	// Metric 3: FCM (Matches filter, Normal latency/error)
	proc.ProcessMetric(models.RawMetric{
		API:       "FCM",
		Timestamp: time.Now(),
		LatencyMs: 100,
		ErrorRate: 0.00,
	})

	// Wait for processing to propagate
	time.Sleep(50 * time.Millisecond)

	// 4. Cancel the context to stop the streaming RPC gracefully
	cancel()

	err := <-errCh
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error on shutdown, got %v", err)
	}

	// 5. Verify results
	if len(stream.messages) != 2 {
		t.Fatalf("expected 2 messages due to API filter, got %d", len(stream.messages))
	}

	// Check First Message (CRITICAL)
	sig1 := stream.messages[0]
	if sig1.Api != "FCM" {
		t.Errorf("expected API FCM, got %s", sig1.Api)
	}
	if sig1.Level != string(models.LevelCritical) {
		t.Errorf("expected level %s, got %s", models.LevelCritical, sig1.Level)
	}

	// Check Second Message (INFO)
	sig2 := stream.messages[1]
	if sig2.Api != "FCM" {
		t.Errorf("expected API FCM, got %s", sig2.Api)
	}
	if sig2.Level != string(models.LevelInfo) {
		t.Errorf("expected level %s, got %s", models.LevelInfo, sig2.Level)
	}
}
