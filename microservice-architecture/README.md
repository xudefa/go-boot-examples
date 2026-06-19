# 微服务架构示例

使用 go-boot 框架构建的微服务架构示例，演示服务拆分、API 网关路由、服务间 HTTP 通信和 IoC 容器集成。

## 架构概述

```
                  ┌─────────────┐
                  │  API Gateway │  :8080
                  └──────┬──────┘
                         │
            ┌────────────┼────────────┐
            │                         │
     /users/*                 /orders/*
            │                         │
            ▼                         ▼
┌───────────────────┐     ┌───────────────────┐
│   User Service    │     │  Order Service    │
│   :8083           │     │   :8084           │
└───────────────────┘     └───────────────────┘
```

## 服务组成

| 服务 | 端口 | 职责 |
|------|------|------|
| API Gateway | 8080 | 统一入口、请求转发 |
| User Service | 8083 | 用户 CRUD |
| Order Service | 8084 | 订单 CRUD、用户验证 |

## 快速开始

需要在 **三个终端窗口** 中分别启动服务：

### 终端 1：启动用户服务

```bash
cd examples/microservice-architecture/user-service
go mod tidy
go run main.go
```

### 终端 2：启动订单服务

```bash
cd examples/microservice-architecture/order-service
go mod tidy
go run main.go
```

### 终端 3：启动 API 网关

```bash
cd examples/microservice-architecture/api-gateway
go mod tidy
go run main.go
```

所有服务启动后，通过 API 网关访问：http://localhost:8080

## 使用示例

### 创建用户

```bash
curl -X POST http://localhost:8080/users/api/users \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","email":"alice@example.com"}'
```

### 查询用户

```bash
curl http://localhost:8080/users/api/users/1
```

### 创建订单

```bash
curl -X POST http://localhost:8080/orders/api/orders \
  -H 'Content-Type: application/json' \
  -d '{"user_id":1,"product_ids":[1,2],"total_amount":99.9}'
```

### 查询订单

```bash
curl http://localhost:8080/orders/api/orders/1
```

### 健康检查

```bash
curl http://localhost:8080/health
curl http://localhost:8083/health
curl http://localhost:8084/health
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.New()` | 每个服务独立的 IoC 容器 |
| `ginadapter.New()` | Gin 与容器集成 |
| HTTP 请求转发 | API 网关代理模式 |
| 服务间通信 | Order Service 调用 User Service |
| 服务拆分 | 独立部署、独立端口 |