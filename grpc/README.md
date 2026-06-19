# gRPC 完整示例

演示 go-boot 框架的 gRPC 服务端与客户端集成，包括一元 RPC 和服务端流式 RPC。

## 目录结构

```
grpc/
├── pb/                  # Protobuf 定义和生成代码
│   ├── hello.proto      # 服务定义
│   ├── hello.pb.go      # 生成的 Protobuf 代码
│   └── hello_grpc.pb.go # 生成的 gRPC 代码
├── grpc-server/         # gRPC 服务端示例
├── grpc-client/         # gRPC 客户端示例
├── grpc-server-health/  # 带健康检查的服务端
└── interceptor-demo/    # 拦截器和追踪示例
```

## 服务定义

```protobuf
service HelloService {
  rpc SayHello(HelloRequest) returns (HelloResponse);
  rpc SayHelloStream(HelloRequest) returns (stream HelloResponse);
}
```

| RPC 方法 | 类型 | 说明 |
|----------|------|------|
| `SayHello` | 一元 RPC | 接收名字返回问候语 |
| `SayHelloStream` | 服务端流式 RPC | 返回 5 条顺序问候消息 |

## 快速开始

### 1. 安装依赖

```bash
# 安装 protoc（macOS）
brew install protobuf

# 安装 protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 安装 protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2. 生成代码

```bash
cd examples/grpc
protoc --proto_dir=pb --go_out=pb --go-grpc_out=pb \
  --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
  pb/hello.proto
```

### 3. 启动服务端

```bash
cd examples/grpc/grpc-server
go run .
```

### 4. 启动客户端（新终端）

```bash
cd examples/grpc/grpc-client
go run .
```

## 预期输出

**服务端：**
```
=== gRPC Server Example (go-boot) ===

HelloService registered
Starting gRPC server on :50051...
SayHello received: name=go-boot
SayHelloStream received: name=Stream
```

**客户端：**
```
=== gRPC Client Example (go-boot) ===

Connected to gRPC server

--- SayHello (Unary RPC) ---
Response: Hello, go-boot! (from go-boot gRPC)

--- SayHelloStream (Server Streaming RPC) ---
Stream response: Hello Stream #1
Stream response: Hello Stream #2
Stream response: Hello Stream #3
Stream response: Hello Stream #4
Stream response: Hello Stream #5
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| gRPC 服务端 | 高性能 RPC 服务端 |
| gRPC 客户端 | RPC 客户端调用 |
| 一元 RPC | 请求-响应模式 |
| 服务端流式 RPC | 服务端推送多条数据 |
| Protocol Buffers | 序列化协议 |