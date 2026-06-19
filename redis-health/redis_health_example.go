// Package main 演示 go-boot 运维监控模块中的 Redis 健康检查功能
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/xudefa/go-boot/actuator"
	"github.com/xudefa/go-boot/boot"
	"github.com/xudefa/go-boot/health"
)

func main() {
	// 创建应用上下文
	app, err := boot.NewApplication(
		boot.WithAppName("redis-health-example"),
	)
	if err != nil {
		log.Fatal("Failed to create application:", err)
	}

	// 启动应用
	if err := app.Start(); err != nil {
		log.Fatal("Failed to start application:", err)
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建 Redis 健康指示器 - 模拟连接检查
	redisIndicator := actuator.NewRedisHealthIndicator(func(ctx context.Context) error {
		// 模拟 Redis 连接检查
		// 在实际应用中，这里会执行真实的 Redis ping 或其他检查操作
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond): // 模拟网络延迟
			// 模拟检查成功
			return nil
		}
	})

	// 创建健康聚合器并添加 Redis 健康指示器
	aggregator := health.NewAggregator()
	aggregator.AddIndicator(redisIndicator)

	fmt.Println("开始健康检查...")

	// 执行健康检查
	result := aggregator.Aggregate(ctx)

	fmt.Printf("Redis 健康状态: %s\n", result.Status)
	if result.Details != nil {
		fmt.Printf("详细信息: %+v\n", result.Details)
	}

	// 验证 Redis 健康指示器的名称
	fmt.Printf("Redis 健康指示器名称: %s\n", redisIndicator.Name())

	// 模拟失败的情况
	fmt.Println("\n模拟 Redis 连接失败的情况...")
	failedRedisIndicator := actuator.NewRedisHealthIndicator(func(ctx context.Context) error {
		return fmt.Errorf("模拟 Redis 连接失败")
	})

	failedAggregator := health.NewAggregator()
	failedAggregator.AddIndicator(failedRedisIndicator)

	failedResult := failedAggregator.Aggregate(ctx)
	fmt.Printf("失败情况下的健康状态: %s\n", failedResult.Status)
	if failedResult.Details != nil {
		fmt.Printf("失败详细信息: %+v\n", failedResult.Details)
	}

	// 停止应用
	app.Stop()
}
