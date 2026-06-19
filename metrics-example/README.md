# 指标收集示例

演示 go-boot 框架的指标收集功能，展示如何使用 Counter（计数器）和 Gauge（仪表盘）收集和展示运行时指标数据。

## 功能特性

- **指标注册表**：`metrics.SimpleRegistry` 统一管理所有指标
- **Counter（计数器）**：单调递增的计数器，适用于请求次数、错误次数等
- **Gauge（仪表盘）**：可增可减的仪表盘，适用于内存使用量、连接数等
- **标签系统**：每个指标可携带标签（tag）进行多维度的数据分类
- **独立使用**：Counter 和 Gauge 可脱离注册表独立使用

## 快速开始

```bash
cd examples/metrics-example
go mod tidy
go run main.go
```

## 预期输出

```
=== Metrics Example ===

Metrics registry created
http.requests.total counter: 5

memory.usage gauge: 55.7
memory.usage gauge after set: 32.1
memory.usage(stack) gauge: 12.3
```

## 代码结构

| 组件 | 说明 |
|------|------|
| `SimpleRegistry` | 指标注册表，管理所有指标的注册和收集 |
| `SimpleCounter` | 计数器实现，支持 `Inc()` 和 `Add(n)` 操作 |
| `SimpleGauge` | 仪表盘实现，支持 `Set(n)`、`Add(n)` 和 `Dec()` 操作 |
| `Metric` | 指标数据结构，包含名称、值和标签 |

## 使用示例

### 创建计数器

```go
counter := metrics.NewCounter("http.requests.total")
counter.Inc()
counter.Add(5)
```

### 创建仪表盘

```go
gauge := metrics.NewGauge("memory.usage")
gauge.Set(55.7)
gauge.Dec()
```

### 使用注册表

```go
registry := metrics.NewSimpleRegistry()
registry.Register(counter)
registry.Register(gauge)
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `SimpleRegistry` | 指标注册与管理 |
| `SimpleCounter` | 单调递增计数 |
| `SimpleGauge` | 可增可减的仪表盘 |
| 标签系统 | 多维度数据分类 |