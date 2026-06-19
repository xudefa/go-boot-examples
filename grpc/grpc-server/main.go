// Package main 演示 go-boot 的 gRPC 服务端启动
//
// 本示例展示如何:
//   - 使用 grpc/server 包创建服务端
//   - 实现一元 RPC 方法
//   - 实现服务端流式 RPC 方法
//
// 使用方式:
//
//	cd examples/grpc/grpc-server
//	go run .
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "example-grpc/pb"

	"github.com/xudefa/go-boot-grpc/server"
	"google.golang.org/grpc"
)

// HelloService 是 HelloServiceServer 接口的具体实现，处理客户端的 RPC 请求
type HelloService struct {
	pb.UnimplementedHelloServiceServer
}

// SayHello 处理一元 RPC 请求，接收名字并返回格式化的问候消息
func (s *HelloService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("SayHello received: name=%s", req.Name)
	return &pb.HelloResponse{Message: fmt.Sprintf("Hello, %s! (from go-boot gRPC)", req.Name)}, nil
}

// SayHelloStream 是服务端流式 RPC：依次发送 5 条问候消息，每条间隔 500 毫秒
func (s *HelloService) SayHelloStream(req *pb.HelloRequest, stream grpc.ServerStreamingServer[pb.HelloResponse]) error {
	log.Printf("SayHelloStream received: name=%s", req.Name)
	for i := 1; i <= 5; i++ {
		resp := &pb.HelloResponse{Message: fmt.Sprintf("Hello %s #%d", req.Name, i)}
		if err := stream.Send(resp); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func main() {
	fmt.Println("=== gRPC Server Example (go-boot) ===")
	fmt.Println()

	// 创建 gRPC 服务端实例，监听 50051 端口，使用标准库日志记录器
	srv := server.New(
		server.WithAddress(":50051"),
		server.WithLogger(log.Default()),
	)

	// 将 HelloService 处理器注册到 gRPC 服务端，客户端即可通过服务名调用
	pb.RegisterHelloServiceServer(srv.GRPCServer(), &HelloService{})

	fmt.Println("HelloService registered")
	fmt.Println("Starting gRPC server on :50051...")
	fmt.Println("Press Ctrl+C to stop")

	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
