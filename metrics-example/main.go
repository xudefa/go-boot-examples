// Package main 演示 go-boot 的指标收集
//
// 本示例展示如何:
//   - 创建指标注册表
//   - 使用 Counter（计数器）
//   - 使用 Gauge（仪表盘）
//   - 使用 Histogram（直方图）
//   - 收集所有指标
//   - 使用独立的 Counter/Gauge/Histogram
//   - 使用导出器
//
// 使用方式:
//
//	cd examples/metrics-example
//	go run .
package main

import (
	"fmt"

	"github.com/xudefa/go-boot/metrics"
)

func main() {
	fmt.Println("=== Metrics Example ===")
	fmt.Println()

	// 创建指标注册表，统一管理所有 Counter、Gauge 和 Histogram
	// 注册表支持按名称和标签维度区分不同的指标实例
	registry := metrics.NewSimpleRegistry()
	fmt.Println("Metrics registry created")

	// 创建 Counter（计数器），用于统计 HTTP 请求总数
	// 标签 "method" 和 "path" 用于区分不同的请求维度的指标
	fmt.Println("\n=== Counter 示例 ===")
	counter := registry.Counter("http.requests.total", "method", "GET", "path", "/api/users")
	counter.Inc()
	counter.Inc()
	counter.Add(3)
	fmt.Printf("http.requests.total (GET /api/users): %.0f\n", counter.Value())

	// 创建不同标签的 Counter
	counterPost := registry.Counter("http.requests.total", "method", "POST", "path", "/api/users")
	counterPost.Inc()
	counterPost.Add(2)
	fmt.Printf("http.requests.total (POST /api/users): %.0f\n", counterPost.Value())

	// 创建 Gauge（仪表盘），用于监控堆内存使用量
	// Gauge 的值可以增加也可以减少，适合表示当前快照值
	fmt.Println("\n=== Gauge 示例 ===")
	gauge := registry.Gauge("memory.usage", "type", "heap")
	gauge.Set(45.5)
	gauge.Add(10.2)
	fmt.Printf("memory.usage (heap): %.1f MB\n", gauge.Value())
	gauge.Set(32.1)
	fmt.Printf("memory.usage (heap) after set: %.1f MB\n", gauge.Value())

	// 创建同名的 Gauge，但通过 "type" 标签区分不同维度
	gaugeStack := registry.Gauge("memory.usage", "type", "stack")
	gaugeStack.Set(12.3)
	fmt.Printf("memory.usage (stack): %.1f MB\n", gaugeStack.Value())

	// 模拟内存波动
	gauge.Add(5.0)
	fmt.Printf("memory.usage (heap) after add: %.1f MB\n", gauge.Value())
	gauge.Add(-3.5)
	fmt.Printf("memory.usage (heap) after subtract: %.1f MB\n", gauge.Value())

	// 创建 Histogram（直方图），用于记录请求延迟分布
	fmt.Println("\n=== Histogram 示例 ===")
	histogram := registry.Histogram("request.duration", "service", "api", "endpoint", "/users")

	// 记录一些延迟样本
	histogram.Record(15.2)
	histogram.Record(23.5)
	histogram.Record(18.7)
	histogram.Record(45.3)
	histogram.Record(9.8)

	fmt.Printf("request.duration - Count: %d, Sum: %.2f, Avg: %.2f ms\n",
		histogram.Count(), histogram.Sum(), histogram.Sum()/float64(histogram.Count()))

	// 使用 RecordWithLabels 添加额外标签
	histogram.RecordWithLabels(32.1, map[string]string{"status_code": "200"})
	fmt.Printf("After RecordWithLabels - Count: %d, Sum: %.2f\n", histogram.Count(), histogram.Sum())

	// 创建独立的指标（不使用注册表）
	fmt.Println("\n=== 独立指标示例 ===")

	// 独立 Counter
	c := metrics.NewSimpleCounter()
	c.Inc()
	c.Add(5)
	fmt.Printf("Standalone counter: %.0f\n", c.Value())
	c.Reset()
	fmt.Printf("Standalone counter after reset: %.0f\n", c.Value())

	// 独立 Gauge
	g := metrics.NewSimpleGauge()
	g.Set(100)
	g.Add(25)
	g.Set(50)
	fmt.Printf("Standalone gauge: %.0f\n", g.Value())

	// 独立 Histogram
	h := metrics.NewSimpleHistogram("latency", map[string]string{"type": "api"})
	h.Record(50.0)
	h.Record(75.5)
	h.Record(25.3)
	fmt.Printf("Standalone histogram - Count: %d, Sum: %.2f\n", h.Count(), h.Sum())

	// 收集所有指标
	fmt.Println("\n=== 收集所有指标 ===")
	allMetrics := registry.Collect()
	fmt.Printf("Total metrics collected: %d\n", len(allMetrics))
	for _, m := range allMetrics {
		fmt.Printf("  Metric: name=%s, type=%s, value=%.2f, tags=%v",
			m.Name, m.Type, m.Value, m.Tags)
		if m.Type == "histogram" {
			fmt.Printf(", count=%d, sum=%.2f", m.Count, m.Sum)
		}
		fmt.Println()
	}

	// 使用导出器
	fmt.Println("\n=== 导出器示例 ===")
	consoleExporter := metrics.NewConsoleExporter()
	registry.RegisterExporter(consoleExporter)
	fmt.Println("ConsoleExporter registered")

	err := registry.Export()
	if err != nil {
		fmt.Printf("Export failed: %v\n", err)
	} else {
		fmt.Println("Metrics exported successfully")
	}

	// 重置所有指标
	fmt.Println("\n=== 重置指标 ===")
	registry.Reset()
	resetMetrics := registry.Collect()
	fmt.Printf("Metrics after reset: %d\n", len(resetMetrics))
	for _, m := range resetMetrics {
		fmt.Printf("  %s: %.2f\n", m.Name, m.Value)
	}

	fmt.Println("\nmetrics-example completed successfully!")
}
