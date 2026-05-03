package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/paddyzab/gcc-data-poller/internal/api"
	"github.com/paddyzab/gcc-data-poller/internal/config"
	"github.com/paddyzab/gcc-data-poller/internal/poller"
	"github.com/paddyzab/gcc-data-poller/internal/processor"
)

func main() {
	log.Println("Starting gcc-data-poller...")

	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Initialize Signal Processor
	sigProcessor := processor.NewSignalProcessor()

	// 3. Initialize Data Poller (Firebase & Identity Toolkit)
	dataPoller := poller.NewDataPoller(cfg, sigProcessor)

	// 4. Initialize gRPC Streaming API Server
	grpcServer := api.NewStreamServer(cfg, sigProcessor)

	// 5. Start Poller
	go func() {
		log.Println("Starting data poller...")
		if err := dataPoller.Start(ctx); err != nil {
			log.Fatalf("Poller encountered an error: %v", err)
		}
	}()

	// 6. Start gRPC Server
	go func() {
		log.Println("Starting gRPC streaming server...")
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("gRPC server encountered an error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gracefully...")
	grpcServer.Stop()
	log.Println("Shutdown complete.")
}
