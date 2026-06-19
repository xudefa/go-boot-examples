package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/xudefa/go-boot-grpc/server"
)

func main() {
	log.Println("=== gRPC Server with Health Check Example (go-boot) ===")

	srv := server.New(
		server.WithAddress(":50051"),
		server.WithHealthCheck(),
	)

	log.Println("Starting gRPC server on :50051 with health check...")
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	healthClient := grpc_health_v1.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "grpc"})
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}

	log.Printf("Health status: %v", resp.Status)

	log.Println("Press Ctrl+C to stop")
	select {}
}
