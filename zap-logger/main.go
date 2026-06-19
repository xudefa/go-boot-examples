// Package main 演示 go-boot 与 Zap - 日志适配器
//
// 本示例展示如何:
//   - 创建 Zap 日志适配器
//   - 使用各种日志级别
//   - 使用构建器模式
//
// 使用方式:
//
//	cd examples/zap-logger
//	go run .
package main

import (
	"context"
	"fmt"

	"github.com/xudefa/go-boot/log"
	"github.com/xudefa/go-boot-zap"
)

func main() {
	fmt.Println("=== Zap Logger Example ===")
	fmt.Println()

	ctx := context.Background()

	// NewZapAdapter 创建 Zap 日志适配器，将 Zap 的 Logger 包装为 go-boot 的 Logger 接口
	// WithZapLevel 设置日志级别，WithZapFormat 选择输出格式（console 或 json）
	adapter := zap.NewZapAdapter(
		zap.WithZapLevel(log.DebugLevel),
		zap.WithZapFormat("console"),
	)
	fmt.Println("Zap logger created")

	// 使用 Zap 适配器输出不同级别的日志，API 与标准 Logger 完全一致
	adapter.Debug(ctx, "This is a debug message", log.KeyValue{Key: "module", Value: "example"})
	adapter.Info(ctx, "This is an info message", log.KeyValue{Key: "user", Value: "alice"})
	adapter.Warn(ctx, "This is a warning", log.KeyValue{Key: "code", Value: 123})
	adapter.Error(ctx, "This is an error", log.KeyValue{Key: "error", Value: "something went wrong"})

	// With 派生带有 requestId 上下文的日志器，后续日志自动携带该字段
	logger2 := adapter.With(ctx, log.KeyValue{Key: "requestId", Value: "req-001"})
	logger2.Info(ctx, "Request processed", log.KeyValue{Key: "duration", Value: "45ms"})

	// 函数式选项模式创建 Zap 适配器
	// WithZapAddCaller(true) 可在日志中包含调用者文件名和行号
	fmt.Println("\n=== Functional options mode ===")
	builderLogger := zap.NewZapAdapter(
		zap.WithZapLevel(log.InfoLevel),
		zap.WithZapFormat("json"),
		zap.WithZapAddCaller(true),
	)
	builderLogger.Info(ctx, "Built with functional options", log.KeyValue{Key: "method", Value: "options"})

	adapter.Sync()
	fmt.Println("\nzap-logger example completed successfully!")
}
