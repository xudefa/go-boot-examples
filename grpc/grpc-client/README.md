# gRPC 客户端示例

演示 go-boot 的 gRPC 客户端调用。

## 功能特性

- 使用 `grpc/client` 包连接 gRPC 服务端
- 一元 RPC 调用：`SayHello` — 发送名字获取问候语
- 服务端流式 RPC 调用：`SayHelloStream` — 接收多条流式响应

## 快速开始

先启动 gRPC 服务端，再启动客户端：

```bash
# 终端 1: 启动服务端
cd examples/grpc/grpc-server && go run .

# 终端 2: 启动客户端
cd examples/grpc/grpc-client && go run .
```

## 代码结构

```go
// 创建 gRPC 客户端
cli := client.New(
    client.WithAddress("localhost:50051"),
    client.WithTimeout(5*time.Second),
)

// 一元 RPC 调用
resp, err := cli.SayHello(ctx, &pb.HelloRequest{Name: "go-boot"})

// 流式 RPC 调用
stream, err := cli.SayHelloStream(ctx, &pb.HelloRequest{Name: "Stream"})
for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break
    }
    fmt.Println("Stream response:", resp.Message)
}
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `client.New()` | 创建 gRPC 客户端 |
| 一元 RPC 调用 | 请求-响应模式 |
| 流式 RPC 调用 | 接收服务端推送数据 |
| 超时控制 | 请求超时设置 |

## 依赖

- `github.com/xudefa/go-boot/grpc/client` — gRPC 客户端封装
- `google.golang.org/grpc` — gRPC 框架