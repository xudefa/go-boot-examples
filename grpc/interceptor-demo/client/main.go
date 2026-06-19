package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "example-grpc/interceptor-demo/pb"

	"github.com/xudefa/go-boot-grpc/client"
	"github.com/xudefa/go-boot/tracing"
)

func main() {
	fmt.Println("=== gRPC Interceptor & Tracing Client Demo ===")
	fmt.Println()

	tracer := tracing.GetTracer("grpc-client-demo")

	cli := client.New(
		client.WithAddress("localhost:50052"),
		client.WithTimeout(5*time.Second),
		client.WithTracing(tracer),
	)

	if err := cli.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer cli.Close()

	fmt.Println("✓ Client connected with tracing enabled")
	fmt.Println()

	demoClient := pb.NewDemoServiceClient(cli.Conn())
	pingClient := pb.NewPingServiceClient(cli.Conn())

	fmt.Println("--- Testing DemoService.Echo (has service-level interceptor) ---")
	time.Sleep(100 * time.Millisecond)

	echoResp, err := demoClient.Echo(context.Background(), &pb.EchoRequest{Message: "Hello World"})
	if err != nil {
		log.Fatalf("DemoService.Echo failed: %v", err)
	}
	fmt.Printf("Response: %s\n", echoResp.Message)
	fmt.Println()

	fmt.Println("--- Testing PingService.Ping (no service-level interceptor) ---")
	time.Sleep(100 * time.Millisecond)

	pingResp, err := pingClient.Ping(context.Background(), &pb.PingRequest{})
	if err != nil {
		log.Fatalf("PingService.Ping failed: %v", err)
	}
	fmt.Printf("Response: %s\n", pingResp.Message)
	fmt.Println()

	fmt.Println("=== Demo Complete ===")
	fmt.Println()
	fmt.Println("Key Features Demonstrated:")
	fmt.Println("1. Unified Interceptor Registry - RegisterGlobal() / RegisterService()")
	fmt.Println("2. Service-level vs Global Interceptors - Service-specific filtering")
	fmt.Println("3. Automatic Tracing Integration - WithTracing() option")
	fmt.Println("4. Client-side Interceptors - ClientInterceptorRegistry")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("- Global interceptors apply to all services (Auth, Logging)")
	fmt.Println("- Service interceptors only apply to matching services (DemoService)")
	fmt.Println("- Tracing automatically enabled when WithTracing() is used")
	fmt.Println("- Interceptors are chained in registration order")
}
