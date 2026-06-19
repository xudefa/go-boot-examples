# gRPC 服务端健康检查示例

演示 go-boot 的 gRPC 服务端健康检查功能。

## 功能特性

- 使用 `WithHealthCheck()` 选项启用 `grpc.health.v1.Health` 服务
- 客户端可调用健康检查端点获取服务状态

## 快速开始

```bash
cd examples/grpc/grpc-server-health
go mod tidy
go run .
```

## 预期输出

```
=== gRPC Server with Health Check Example (go-boot) ===
Starting gRPC server on :50051 with health check...
gRPC server listening on :50051
Health status: SERVING
Press Ctrl+C to stop
```

## 代码结构

```go
// 创建带健康检查的 gRPC 服务端
srv := server.New(
    server.WithAddress(":50051"),
    server.WithHealthCheck(),
)
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `WithHealthCheck()` | 启用健康检查服务 |
| `grpc.health.v1` | gRPC 标准健康检查协议 |
| 服务状态 | SERVING/NOT_SERVING/UNKNOWN |

## 依赖

- `github.com/xudefa/go-boot/grpc/server` — gRPC 服务端封装
- `google.golang.org/grpc/health/grpc_health_v1` — gRPC 健康检查协议