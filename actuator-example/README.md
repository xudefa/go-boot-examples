# Actuator 运维端点示例

演示 go-boot 框架的 Actuator 运维端点功能，展示如何暴露健康检查、指标和运行时信息的 HTTP 接口。

## 功能特性

- **健康检查端点** (`/actuator/health`)：聚合展示各组件的健康状态
- **指标端点** (`/actuator/metrics`)：展示运行时指标数据
- **环境信息端点** (`/actuator/env`)：展示应用配置和环境信息
- **Bean 信息端点** (`/actuator/beans`)：展示 IoC 容器中的 Bean 列表
- **自定义健康指标**：支持自定义组件健康检查逻辑

## 快速开始

```bash
cd examples/actuator-example
go mod tidy
go run main.go
```

服务启动后可访问以下端点：

| 端点 | 方法 | 说明 |
|------|------|------|
| `/actuator/health` | GET | 聚合健康状态 |
| `/actuator/metrics` | GET | 运行时指标 |
| `/actuator/env` | GET | 环境配置信息 |
| `/actuator/beans` | GET | Bean 列表 |

## 代码结构

| 组件 | 说明 |
|------|------|
| `MyAppIndicator` | 自定义应用健康指标，返回版本和运行时间 |
| `NewApplication` | 创建应用实例，配置应用名称和版本 |
| `NewDatabaseHealthIndicator` | 数据库健康检查（模拟） |
| `NewRedisHealthIndicator` | Redis 健康检查（模拟） |
| `RegisterRoutes` | 注册所有 Actuator HTTP 路由 |

## 使用示例

### 查看健康状态

```bash
curl http://localhost:9090/actuator/health
```

### 查看指标信息

```bash
curl http://localhost:9090/actuator/metrics
```

### 查看环境配置

```bash
curl http://localhost:9090/actuator/env
```

### 查看 Bean 列表

```bash
curl http://localhost:9090/actuator/beans
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| Actuator 端点 | 运维监控接口 |
| 健康指标聚合 | 多组件健康状态汇总 |
| 自定义 Indicator | 扩展健康检查逻辑 |
| HTTP 路由注册 | Actuator 端点暴露 |