// 示例应用程序：展示 go-boot 运维监控模块的使用
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xudefa/go-boot/boot"
	"github.com/xudefa/go-boot/health"
	"github.com/xudefa/go-boot/metrics"
	"github.com/xudefa/go-boot-prometheus"
	"github.com/xudefa/go-boot/tracing"
)

// databaseIndicator 实现健康检查指示器接口
type databaseIndicator struct {
	name string
}

func (d *databaseIndicator) Name() string {
	return d.name
}

func (d *databaseIndicator) Health(ctx context.Context) health.Health {
	// 模拟数据库检查
	return health.Health{
		Status:  health.StatusUp,
		Details: map[string]interface{}{"database": "connected", "version": "PostgreSQL 13.4"},
	}
}

func main() {
	var metricRegistry metrics.MeterRegistry
	var prometheusExporter *prometheus.Exporter

	// 创建应用上下文
	app, err := boot.NewApplication(
		boot.WithAppName("monitoring-example"),
	)
	if err != nil {
		log.Fatal("Failed to create application:", err)
	}

	// 启动应用
	if err := app.Start(); err != nil {
		log.Fatal("Failed to start application:", err)
	}

	// 获取指标注册表
	obj, err := app.Context().Get("meterRegistry")
	if err != nil {
		log.Printf("Warning: Failed to get metric registry: %v", err)
		// 如果获取不到指标注册表，创建一个新的
		metricRegistry = metrics.NewSimpleRegistry()
	} else {
		var ok bool
		metricRegistry, ok = obj.(metrics.MeterRegistry)
		if !ok {
			log.Printf("Warning: Failed to cast object to metrics.MeterRegistry")
			metricRegistry = metrics.NewSimpleRegistry()
		}
	}

	// 获取 Prometheus 导出器
	obj2, err := app.Context().Get("prometheus.exporter")
	if err != nil {
		log.Printf("Warning: Failed to get prometheus exporter: %v", err)
		// 如果获取不到 Prometheus 导出器，创建一个新的
		prometheusExporter = prometheus.NewExporter(metricRegistry)
		if err := prometheusExporter.Start(); err != nil {
			log.Printf("Warning: Failed to start prometheus exporter: %v", err)
		} else {
			defer prometheusExporter.Stop()
		}
	} else {
		var ok bool
		prometheusExporter, ok = obj2.(*prometheus.Exporter)
		if !ok {
			log.Printf("Warning: Failed to cast object to *prometheus.Exporter")
			prometheusExporter = prometheus.NewExporter(metricRegistry)
			if err := prometheusExporter.Start(); err != nil {
				log.Printf("Warning: Failed to start prometheus exporter: %v", err)
			} else {
				defer prometheusExporter.Stop()
			}
		}
	}

	// 创建追踪器
	tracer := tracing.NewTracer("example-service")

	// 创建自定义健康检查指示器
	dbIndicator := &databaseIndicator{name: "database"}

	// 创建聚合器并添加指示器
	aggregator := health.NewAggregator()
	aggregator.AddIndicator(dbIndicator)

	// 模拟业务操作
	if metricRegistry != nil {
		go func() {
			counter := metricRegistry.Counter("business_operations_total", "type", "process")
			gauge := metricRegistry.Gauge("current_processing_items")

			for i := 0; i < 100; i++ {
				// 记录操作计数
				counter.Inc()

				// 更新当前处理项目数
				gauge.Set(float64(i % 10))

				// 模拟业务操作
				func() {
					_, span := tracer.Start(context.Background(), "business-operation")
					defer span.End()

					span.SetAttribute("operation.id", fmt.Sprintf("op-%d", i))

					// 模拟处理时间
					time.Sleep(100 * time.Millisecond)
				}()

				time.Sleep(500 * time.Millisecond)
			}
		}()
	}

	// 启动 HTTP 服务器，暴露 Prometheus 指标端点
	http.Handle("/metrics", prometheusExporter.Handler())

	// 添加健康检查端点
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		result := aggregator.Aggregate(context.Background())
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"%s","details":%v}`, result.Status, result.Details)
	})

	fmt.Println("应用已启动")
	fmt.Println("- 指标端点: http://localhost:9090/metrics")
	fmt.Println("- 健康检查端点: http://localhost:9090/health")

	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}

	// 停止应用
	app.Stop()
}
