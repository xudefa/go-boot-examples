# gRPC 拦截器与追踪示例

演示 go-boot gRPC 模块的拦截器和追踪集成功能。

## 功能特性

- **统一的拦截器注册机制**：使用 `InterceptorRegistry` 管理所有拦截器
- **服务级别和全局拦截器**：支持全局拦截器和服务特定拦截器
- **自动追踪集成**：通过 `WithTracing()` 选项自动启用分布式追踪
- **客户端拦截器**：客户端也支持拦截器链
- **线程安全**：拦截器注册是线程安全的

## 目录结构

```
interceptor-demo/
├── demo.proto          # Protobuf 定义文件
├── pb/                 # 生成的 Protobuf 代码
├── server/             # 服务端示例
│   └── main.go
├── client/             # 客户端示例
│   └── main.go
└── README.md           # 本文件
```

## 快速开始

### 1. 生成 Protobuf 代码

```bash
cd examples/grpc/interceptor-demo
protoc --go_out=. --go-grpc_out=. --proto_path=. demo.proto
```

### 2. 启动服务端

```bash
cd server
go mod tidy
go run .
```

服务端将在 `:50052` 端口监听。

### 3. 启动客户端（新终端）

```bash
cd client
go mod tidy
go run .
```

## 服务定义

| 服务 | 方法 | 说明 |
|------|------|------|
| `DemoService` | `Echo` | 回显消息 |
| `PingService` | `Ping` | 返回 Pong |

## 拦截器配置

### 全局拦截器（应用于所有服务）

| 拦截器 | 说明 |
|--------|------|
| Auth Interceptor | 模拟认证检查 |
| Logging Interceptor | 记录 RPC 调用耗时和错误 |

### 服务级别拦截器

| 拦截器 | 说明 |
|--------|------|
| DemoService Interceptor | 仅应用于 `DemoService` 服务 |

## 代码结构

### 服务端拦截器注册

```go
srv := server.New(
    server.WithAddress(":50052"),
    server.WithTracing(tracer),
    server.WithGlobalInterceptor(authInterceptor),
    server.WithGlobalInterceptor(loggingInterceptor),
    server.WithServiceInterceptor("/demo.DemoService", demoServiceInterceptor),
)
```

### 客户端拦截器配置

```go
cli := client.New(
    client.WithAddress("localhost:50052"),
    client.WithTimeout(5*time.Second),
    client.WithTracing(tracer),
)
```

## 拦截器执行顺序

拦截器按照注册顺序执行，后注册的拦截器先执行（洋葱模型）。

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| 全局拦截器 | 应用于所有服务 |
| 服务级别拦截器 | 仅应用于特定服务 |
| 分布式追踪 | 自动集成 OpenTelemetry |
| 客户端拦截器 | 客户端请求拦截 |