# gRPC 服务端示例

演示 go-boot 的 gRPC 服务端启动与 RPC 方法实现。

## 功能特性

- 使用 `grpc/server` 包创建 gRPC 服务端
- 一元 RPC：`SayHello` — 接收名字返回问候语
- 服务端流式 RPC：`SayHelloStream` — 返回 5 条顺序问候消息
- Protocol Buffers 协议定义服务接口

## 快速开始

```bash
cd examples/grpc/grpc-server
go mod tidy
go run .
```

服务端监听在 `:50051`。

## 代码结构

```go
// 创建 gRPC 服务端
srv := server.New(server.WithAddress(":50051"))

// 注册服务
pb.RegisterHelloServiceServer(srv.GRPCServer(), &helloServer{})

// 启动服务端
srv.Start()
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `server.New()` | 创建 gRPC 服务端 |
| 服务注册 | 注册 Protobuf 服务 |
| 一元 RPC | 请求-响应模式 |
| 服务端流式 RPC | 服务端推送数据 |

## 依赖

- `github.com/xudefa/go-boot/grpc/server` — gRPC 服务端封装
- `google.golang.org/grpc` — gRPC 框架
- `google.golang.org/protobuf` — Protobuf 序列化