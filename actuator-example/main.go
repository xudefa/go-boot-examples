// Package main 演示 go-boot 的 Actuator（运维端点）
//
// 本示例展示如何:
//   - 创建自定义健康指标
//   - 使用数据库和 Redis 健康指标
//   - 注册 Actuator HTTP 路由（健康检查、指标、环境信息）
//
// 使用方式:
//
//	cd examples/actuator-example
//	go run .
//	# 访问 http://localhost:9090/actuator/health
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xudefa/go-boot/actuator"
	"github.com/xudefa/go-boot/boot"
	"github.com/xudefa/go-boot/health"
)

// MyAppIndicator 自定义应用健康指标
// 实现 health.Indicator 接口，返回应用自身的健康状态信息
type MyAppIndicator struct{}

// Name 返回该指标的标识名称，在健康检查结果中作为 key 展示
func (m *MyAppIndicator) Name() string {
	return "myapp"
}

// Health 执行健康检查逻辑，返回应用版本和运行时间等信息
func (m *MyAppIndicator) Health(ctx context.Context) health.Health {
	return health.Health{
		Status: health.StatusUp,
		Details: map[string]any{
			"version": "1.0.0",
			"uptime":  "5m",
		},
	}
}

func main() {
	fmt.Println("=== Actuator Example ===")
	fmt.Println()

	// 使用 boot.NewApplication 创建应用实例
	// 配置应用名称和版本，应用上下文会提供环境信息给 Actuator
	app, err := boot.NewApplication(
		boot.WithAppName("actuator-demo"),
		boot.WithVersion("1.0.0"),
	)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}
	fmt.Println("Application created")

	// 创建 Actuator 实例，传入应用上下文
	// Actuator 负责管理运维端点的注册和路由分发
	act := actuator.New(app.Context())
	fmt.Println("Actuator created")

	// 创建健康聚合器并注册自定义的健康指标
	agg := health.NewAggregator()
	agg.AddIndicator(&MyAppIndicator{})
	fmt.Println("Health aggregator configured with myapp indicator")

	// 创建数据库健康指标，使用回调函数模拟健康检查逻辑
	// 回调返回 nil 表示数据库连接正常
	dbHealth := actuator.NewDatabaseHealthIndicator(func(ctx context.Context) error {
		return nil
	})
	agg.AddIndicator(dbHealth)
	fmt.Println("Database health indicator added")

	// 创建 Redis 健康指标，同样使用回调模拟
	redisHealth := actuator.NewRedisHealthIndicator(func(ctx context.Context) error {
		return nil
	})
	agg.AddIndicator(redisHealth)
	fmt.Println("Redis health indicator added")

	// 将健康聚合器设置到 Actuator 中，使其可以通过 HTTP 暴露健康数据
	act.SetHealthAggregator(agg)
	// 获取 Actuator 内部的指标注册表，用于收集和暴露运行时指标
	registry := act.MetricsRegistry()
	fmt.Printf("Metrics registry obtained: %v\n", registry)

	// 注册 Actuator HTTP 路由到默认的 ServeMux
	// 注册后可通过 HTTP 访问各运维端点
	mux := http.NewServeMux()
	act.RegisterRoutes(mux, actuator.DefaultRouteConfig())
	fmt.Println("Actuator routes registered:")
	fmt.Println("  GET /actuator/health")
	fmt.Println("  GET /actuator/metrics")
	fmt.Println("  GET /actuator/env")
	fmt.Println("  GET /actuator/beans")

	// 启动 HTTP 服务器，监听 9090 端口
	// 在 goroutine 中异步启动，避免阻塞主流程
	server := &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}
	go func() {
		fmt.Println("\nStarting actuator server on :9090...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// 等待服务器启动后，使用带超时的上下文优雅关闭
	time.Sleep(100 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	fmt.Println("\nServer stopped")

	fmt.Println("\nactuator-example completed successfully!")
}
