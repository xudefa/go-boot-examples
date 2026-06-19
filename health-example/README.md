# 健康检查示例

演示 go-boot 框架的健康检查功能，包括自定义健康指标和健康状态聚合。

## 功能特性

- **自定义健康指标**：实现 `health.Indicator` 接口创建自定义组件健康检查
- **健康聚合器**：使用 `health.Aggregator` 聚合多个健康指标的状态
- **健康状态枚举**：五种健康状态

## 健康状态

| 状态 | 说明 |
|------|------|
| `StatusUp` | 正常运行 |
| `StatusDown` | 不可用 |
| `StatusDegraded` | 性能降级 |
| `StatusOutage` | 服务中断 |
| `StatusUnknown` | 状态未知 |

## 快速开始

```bash
cd examples/health-example
go mod tidy
go run main.go
```

## 预期输出

```
=== Health Indicator Example ===

Health aggregator created
Indicators added: database, redis
Total indicators: 2

Aggregated health status: UP
Details: map[database:map[connected:true latency:5ms] redis:map[connected:true latency:2ms]]
```

## 代码结构

| 组件 | 说明 |
|------|------|
| `DatabaseHealthIndicator` | 模拟数据库健康检查 |
| `RedisHealthIndicator` | 模拟 Redis 健康检查 |
| `health.Aggregator` | 聚合所有注册指标的健康状态 |

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `health.Indicator` 接口 | 自定义健康检查实现 |
| `health.Aggregator` | 多指标健康聚合 |
| 健康状态枚举 | 细粒度的健康状态 |