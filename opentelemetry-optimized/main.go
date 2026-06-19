// Package main 演示 go-boot 的 OpenTelemetry 优化集成
//
// 本示例展示如何:
//   - 配置 OpenTelemetry 追踪
//   - 通过依赖注入注入 Tracer
//   - 创建父子 Span 关系
//   - 记录 Span 属性和事件
//   - 设置 Span 状态和错误
//
// 使用方式:
//
//	cd examples/opentelemetry-optimized
//	go run .
//
// 需要先启动 OpenTelemetry Collector（可选）:
//
//	docker run -p 4317:4317 otel/opentelemetry-collector-contrib:latest
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/xudefa/go-boot-opentelemetry"
	"github.com/xudefa/go-boot/boot"
	"github.com/xudefa/go-boot/core"
	"github.com/xudefa/go-boot/environment"
	"github.com/xudefa/go-boot/tracing"
)

type DemoService struct {
	Tracer tracing.Tracer `inject:"tracer"`
}

func (s *DemoService) RunDemo(ctx context.Context) error {
	if s.Tracer == nil {
		return fmt.Errorf("tracer is nil - injection failed")
	}

	fmt.Println("测试基本追踪...")
	_, span := s.Tracer.Start(ctx, "root-operation",
		tracing.WithSpanKind(tracing.SpanKindInternal),
		tracing.WithAttribute("operation", "demo"),
		tracing.WithAttribute("count", 3),
		tracing.WithAttribute("enabled", true),
	)
	defer span.End()

	for i := 0; i < 3; i++ {
		_, childSpan := s.Tracer.Start(ctx, fmt.Sprintf("sub-operation-%d", i),
			tracing.WithSpanKind(tracing.SpanKindInternal),
			tracing.WithAttribute("iteration", i),
			tracing.WithAttribute("duration_ms", 10*i),
		)
		time.Sleep(10 * time.Millisecond)
		childSpan.AddEvent("work-done", tracing.WithEventAttribute("status", "success"))
		childSpan.SetStatus(tracing.SpanStatusOK)
		childSpan.End()
	}

	span.SetAttribute("total_operations", 3)
	span.SetStatus(tracing.SpanStatusOK)

	fmt.Println("\n测试错误处理...")
	_, errSpan := s.Tracer.Start(ctx, "error-operation",
		tracing.WithSpanKind(tracing.SpanKindInternal),
	)
	errSpan.AddEvent("attempt")
	errSpan.SetError(fmt.Errorf("simulated error"))
	errSpan.SetStatus(tracing.SpanStatusError)
	errSpan.End()

	fmt.Println("\n测试不同类型属性...")
	_, attrSpan := s.Tracer.Start(ctx, "attribute-test",
		tracing.WithSpanKind(tracing.SpanKindInternal),
	)
	attrSpan.SetAttribute("string_value", "hello")
	attrSpan.SetAttribute("int_value", 42)
	attrSpan.SetAttribute("float_value", 3.14)
	attrSpan.SetAttribute("bool_value", true)
	attrSpan.SetStatus(tracing.SpanStatusOK)
	attrSpan.End()

	fmt.Println("\n追踪上下文信息:")
	sc := span.SpanContext()
	fmt.Printf("  TraceID: %s\n", sc.TraceID)
	fmt.Printf("  SpanID: %s\n", sc.SpanID)

	return nil
}

func main() {
	fmt.Println("=== OpenTelemetry 优化示例 ===")
	fmt.Println()

	app, err := boot.NewApplication(
		boot.WithAppName("opentelemetry-optimized-demo"),
		boot.WithVersion("1.0.0"),
		boot.WithProfiles("dev"),
	)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	env := app.Environment()
	env.AddPropertySource(environment.NewMapPropertySource("tracing-config", environment.PriorityNormal, map[string]any{
		"tracing.enabled":           "true",
		"tracing.provider":          "opentelemetry",
		"tracing.service.name":      "demo-service",
		"tracing.service.version":   "1.0.0",
		"tracing.exporter.type":     "otlpgrpc",
		"tracing.exporter.endpoint": "localhost:4317",
		"tracing.sampling":          1.0,
		"tracing.environment":       "development",
	}))

	container := app.Container()

	fmt.Println("Starting application...")
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
	defer func() {
		if err := app.Stop(); err != nil {
			log.Printf("Failed to stop application: %v", err)
		}
	}()

	fmt.Println("Checking beans in container...")
	beans := container.ListBeans()
	for _, bean := range beans {
		fmt.Printf("  - %s (%s)\n", bean.ID, bean.Type)
	}

	// 尝试从容器获取 TracerProvider，如果失败则使用 NoOp 实现
	var tracer tracing.Tracer
	tracerProviderAny, err := container.Get("tracerProvider")
	if err != nil {
		fmt.Printf("Warning: tracerProvider not found in container (%v), using NoOp tracer\n", err)
		tracer = &tracing.NoopTracer{}
	} else {
		tracerProvider, ok := tracerProviderAny.(tracing.TracerProvider)
		if !ok {
			fmt.Printf("Warning: failed to cast tracerProvider, using NoOp tracer\n")
			tracer = &tracing.NoopTracer{}
		} else {
			tracer = tracerProvider.Tracer("demo-service")
		}
	}

	// 创建 DemoService 并手动注入 Tracer
	demoService := &DemoService{Tracer: tracer}

	container.Register("demoService",
		core.Bean(demoService),
		core.Singleton())

	ctx := context.Background()

	if err := demoService.RunDemo(ctx); err != nil {
		log.Printf("Demo failed: %v", err)
	}

	fmt.Println("\nOpenTelemetry 优化示例完成!")
}
