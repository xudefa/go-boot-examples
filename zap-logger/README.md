# Zap 日志集成示例

演示 go-boot 与 Zap 高性能日志库的集成，支持多级别日志、结构化输出和上下文作用域。

## 功能特性

- **Zap 适配器**：通过 `NewZapAdapter` 创建，桥接 go-boot 日志接口
- **日志级别**：支持 Debug / Info / Warn / Error
- **函数式选项**：通过 `WithZapLevel`、`WithZapFormat`、`WithZapAddCaller` 配置
- **上下文作用域**：使用 `With()` 创建带请求 ID 的日志作用域
- **日志同步**：通过 `Sync()` 确保缓冲区日志写入

## 快速开始

```bash
cd examples/zap-logger
go mod tidy
go run main.go
```

## 使用示例

### 创建 Zap 适配器

```go
logger := zap.NewZapAdapter(
    zap.WithZapLevel(zap.DebugLevel),
    zap.WithZapFormat("json"),
    zap.WithZapAddCaller(true),
)
```

### 日志输出

```go
logger.Info("Application started")
logger.Warn("Low disk space", "disk", "/dev/sda1", "free", "10%")
logger.Error("Failed to connect", "error", err)
```

### 上下文作用域

```go
reqLogger := logger.With("request_id", "abc123")
reqLogger.Info("Processing request")
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| Zap 适配器 | 桥接 go-boot 日志接口 |
| 函数式选项 | 灵活的日志配置 |
| 结构化日志 | JSON 格式输出 |
| 上下文作用域 | 请求级别日志追踪 |