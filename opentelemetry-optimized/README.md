# OpenTelemetry 链路追踪示例

演示 go-boot 框架的 OpenTelemetry 集成，实现分布式追踪、类型安全属性转换和优雅关闭。

## 功能特性

- **增强的资源属性**：服务版本、主机名、进程信息等
- **类型安全的属性转换**：保留原始类型（string、int、float、bool）
- **完善的 SpanKind 支持**：Internal、Server、Client、Producer、Consumer
- **优雅关闭机制**：确保追踪数据正确刷新
- **嵌套追踪**：支持多层级的 Span 嵌套
- **错误处理**：正确的错误记录和状态设置

## 快速开始

```bash
cd examples/opentelemetry-optimized
go mod tidy
go run main.go
```

## 预期输出

```
=== OpenTelemetry 优化示例 ===

测试基本追踪...

测试错误处理...

测试不同类型属性...

追踪上下文信息:
  TraceID: 4bf92f3577b34da6a3ce929d0e0e4736
  SpanID: 00f067aa0ba902b7
```

按 `Ctrl+C` 退出程序，将看到优雅关闭信息。

## 使用示例

### 基本追踪

```go
ctx, span := tracer.Start(ctx, "root-operation",
    tracing.WithSpanKind(tracing.SpanKindInternal),
    tracing.WithAttribute("operation", "demo"),
    tracing.WithAttribute("count", 3),
)
defer span.End()
```

### 错误处理

```go
_, errSpan := tracer.Start(ctx, "error-operation",
    tracing.WithSpanKind(tracing.SpanKindInternal),
)
errSpan.SetError(fmt.Errorf("simulated error"))
errSpan.SetStatus(tracing.SpanStatusError)
errSpan.End()
```

### 嵌套追踪

```go
ctx, childSpan := tracer.Start(ctx, "sub-operation",
    tracing.WithSpanKind(tracing.SpanKindInternal),
    tracing.WithAttribute("iteration", i),
)
defer childSpan.End()
```

## SpanKind 说明

| SpanKind | 说明 |
|----------|------|
| `SpanKindInternal` | 内部操作 |
| `SpanKindServer` | 服务器端接收请求 |
| `SpanKindClient` | 客户端发起请求 |
| `SpanKindProducer` | 消息生产者 |
| `SpanKindConsumer` | 消息消费者 |

## 配置说明

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `tracing.enabled` | `false` | 是否启用追踪 |
| `tracing.service.name` | `go-boot-app` | 服务名称 |
| `tracing.service.version` | `1.0.0` | 服务版本 |
| `tracing.exporter.type` | `otlp` | 导出器类型 |
| `tracing.exporter.endpoint` | `localhost:4317` | OTLP 端点 |
| `tracing.sampling` | `1.0` | 采样率（0-1） |

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| OpenTelemetry 集成 | 分布式链路追踪 |
| SpanKind | 追踪角色标识 |
| 类型安全属性 | 保留原始数据类型 |
| 优雅关闭 | 数据刷新机制 |