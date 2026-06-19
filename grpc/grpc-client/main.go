// Package main 演示 go-boot 的 gRPC 客户端调用
//
// 本示例展示如何:
//   - 使用 grpc/client 包创建客户端连接
//   - 发起一元 RPC 调用
//   - 发起服务端流式 RPC 调用
//
// 使用方式:
//
//	cd examples/grpc/grpc-client
//	go run .
package main

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
	"time"

	pb "example-grpc/pb"

	"github.com/xudefa/go-boot-grpc/client"
	sdklog "github.com/xudefa/go-boot/log"
)

func main() {
	fmt.Println("=== gRPC Client Example (go-boot) ===")
	fmt.Println()

	// 创建 gRPC 客户端实例，配置服务端地址、超时时间和日志记录器
	cli := client.New(
		client.WithAddress("localhost:50051"),
		client.WithTimeout(10*time.Second),
		client.WithLogger(sdklog.NewSlogLogger()),
	)

	// 建立与 gRPC 服务端的连接，连接失败则直接退出
	if err := cli.Connect(); err != nil {
		stdlog.Fatalf("Failed to connect: %v", err)
	}
	defer cli.Close()
	fmt.Println("Connected to gRPC server")

	// 通过已有连接创建 HelloService 的客户端存根，用于发起 RPC 调用
	grpcClient := pb.NewHelloServiceClient(cli.Conn())

	fmt.Println("\n--- SayHello (Unary RPC) ---")
	// 发起一元 RPC 调用：发送 HelloRequest，等待并接收单个 HelloResponse
	resp, err := grpcClient.SayHello(context.Background(), &pb.HelloRequest{Name: "go-boot"})
	if err != nil {
		stdlog.Fatalf("SayHello failed: %v", err)
	}
	fmt.Printf("Response: %s\n", resp.Message)

	fmt.Println("\n--- SayHelloStream (Server Streaming RPC) ---")
	// 发起服务端流式 RPC 调用：接收一个返回流，循环读取直到 io.EOF
	stream, err := grpcClient.SayHelloStream(context.Background(), &pb.HelloRequest{Name: "Stream"})
	if err != nil {
		stdlog.Fatalf("SayHelloStream failed: %v", err)
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			stdlog.Fatalf("Stream receive failed: %v", err)
		}
		fmt.Printf("Stream response: %s\n", resp.Message)
	}

	fmt.Println("\nAll RPCs completed successfully")
}
