package main

import (
	"context"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/paddyzab/gcc-data-poller/pkg/pb"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	target := "localhost:" + port
	log.Printf("Connecting to gRPC streaming server at %s...", target)

	// Set up a connection to the server
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewSignalStreamerClient(conn)

	// Subscribe to signals. E.g. filtering to only receive "FCM" signals
	req := &pb.SubscribeRequest{
		ApiFilter: "FCM",
	}

	log.Printf("Subscribing to signals with filter: '%s'", req.ApiFilter)
	stream, err := client.SubscribeSignals(context.Background(), req)
	if err != nil {
		log.Fatalf("Error subscribing to signals: %v", err)
	}

	for {
		signal, err := stream.Recv()
		if err == io.EOF {
			log.Println("Server closed the stream")
			break
		}
		if err != nil {
			log.Fatalf("Error receiving stream data: %v", err)
		}

		log.Printf("Received Signal | API: %s | Level: %s | Message: %s | Timestamp: %v",
			signal.Api, signal.Level, signal.Message, signal.Timestamp.AsTime())
	}
}
