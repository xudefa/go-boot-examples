package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "example-grpc/interceptor-demo/pb"

	"github.com/xudefa/go-boot-grpc/server"
	"github.com/xudefa/go-boot/tracing"
	"google.golang.org/grpc"
)

type DemoServiceImpl struct {
	pb.UnimplementedDemoServiceServer
}

func (s *DemoServiceImpl) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Printf("DemoService.Echo received: message=%s", req.Message)
	return &pb.EchoResponse{Message: "Echo: " + req.Message}, nil
}

type PingServiceImpl struct {
	pb.UnimplementedPingServiceServer
}

func (s *PingServiceImpl) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	log.Printf("PingService.Ping received")
	return &pb.PingResponse{Message: "Pong"}, nil
}

func main() {
	fmt.Println("=== gRPC Interceptor & Tracing Server Demo ===")
	fmt.Println()

	tracer := tracing.GetTracer("grpc-server-demo")

	authInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Printf("[Auth Interceptor] Checking auth for: %s\n", info.FullMethod)
		return handler(ctx, req)
	}

	loggingInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)
		fmt.Printf("[Logging Interceptor] %s took %v, error: %v\n", info.FullMethod, duration, err)
		return resp, err
	}

	demoServiceInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Printf("[DemoService Interceptor] DemoService-specific logic for: %s\n", info.FullMethod)
		return handler(ctx, req)
	}

	srv := server.New(
		server.WithAddress(":50052"),
		server.WithTracing(tracer),
		server.WithGlobalInterceptor(authInterceptor),
		server.WithGlobalInterceptor(loggingInterceptor),
		server.WithServiceInterceptor("/demo.DemoService", demoServiceInterceptor),
	)

	pb.RegisterDemoServiceServer(srv.GRPCServer(), &DemoServiceImpl{})
	pb.RegisterPingServiceServer(srv.GRPCServer(), &PingServiceImpl{})

	fmt.Println("✓ Server started with tracing enabled")
	fmt.Println("✓ Global interceptors: Auth, Logging")
	fmt.Println("✓ Service-level interceptor: DemoService")
	fmt.Println("✓ Services registered: DemoService, PingService")
	fmt.Println()
	fmt.Println("Listening on :50052...")
	fmt.Println("Press Ctrl+C to stop")

	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
