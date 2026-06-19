// Package main 演示 go-boot 的健康检查
//
// 本示例展示如何:
//   - 创建自定义健康指标
//   - 使用 Aggregator 聚合多个指标
//   - 了解健康状态枚举值
//
// 使用方式:
//
//	cd examples/health-example
//	go run .
package main

import (
	"context"
	"fmt"

	"github.com/xudefa/go-boot/health"
)

// DatabaseHealthIndicator 数据库健康指标
// 实现 health.Indicator 接口，模拟数据库组件的健康检查
type DatabaseHealthIndicator struct{}

// Name 返回该指标的标识名称，用于在聚合结果中区分不同组件
func (d *DatabaseHealthIndicator) Name() string {
	return "database"
}

// Health 执行健康检查逻辑，返回当前组件的健康状态
// 此处模拟数据库连接正常，返回 UP 状态和延迟信息
func (d *DatabaseHealthIndicator) Health(ctx context.Context) health.Health {
	return health.Health{
		Status: health.StatusUp,
		Details: map[string]any{
			"connected": true,
			"latency":   "5ms",
		},
	}
}

// RedisHealthIndicator Redis 健康指标
// 模拟 Redis 缓存组件的健康检查
type RedisHealthIndicator struct{}

// Name 返回该指标的标识名称
func (r *RedisHealthIndicator) Name() string {
	return "redis"
}

// Health 模拟 Redis 健康检查，返回 UP 状态和较低延迟
func (r *RedisHealthIndicator) Health(ctx context.Context) health.Health {
	return health.Health{
		Status: health.StatusUp,
		Details: map[string]any{
			"connected": true,
			"latency":   "2ms",
		},
	}
}

func main() {
	fmt.Println("=== Health Indicator Example ===")
	fmt.Println()

	// 创建健康聚合器，用于统一收集和管理多个组件的健康状态
	aggregator := health.NewAggregator()
	fmt.Println("Health aggregator created")

	// 向聚合器中注册自定义的健康指标
	// 支持同时注册多个不同类型的组件指标
	aggregator.AddIndicator(&DatabaseHealthIndicator{})
	aggregator.AddIndicator(&RedisHealthIndicator{})
	fmt.Println("Indicators added: database, redis")

	// 查询当前已注册的指标总数
	fmt.Printf("Total indicators: %d\n", len(aggregator.Indicators()))

	// 执行聚合操作，遍历所有指标并汇总健康状态
	// 如果所有指标都是 UP，整体状态为 UP；任一指标 DOWN 则整体为 DOWN
	result := aggregator.Aggregate(context.Background())
	fmt.Printf("\nAggregated health status: %s\n", result.Status)
	fmt.Printf("Details: %+v\n", result.Details)
	fmt.Printf("Timestamp: %s\n", result.Timestamp)

	fmt.Println("\n=== Status values ===")
	// 健康状态枚举值列表：
	//   StatusUp      - 组件正常运行
	//   StatusDown    - 组件不可用
	//   StatusDegraded - 组件降级运行
	//   StatusOutage  - 组件完全中断
	//   StatusUnknown - 组件状态未知
	fmt.Printf("  StatusUp:      %s\n", health.StatusUp)
	fmt.Printf("  StatusDown:    %s\n", health.StatusDown)
	fmt.Printf("  StatusDegraded:%s\n", health.StatusDegraded)
	fmt.Printf("  StatusOutage:  %s\n", health.StatusOutage)
	fmt.Printf("  StatusUnknown: %s\n", health.StatusUnknown)

	fmt.Println("\nhealth example completed successfully!")
}
