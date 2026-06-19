# 订单服务

微服务架构中的订单管理服务，负责订单的创建、查询、状态更新，并通过 HTTP 调用用户服务验证用户存在性。

## 端口

`http://localhost:8084`

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /health | 健康检查 |
| GET | /api/orders | 查询所有订单 |
| GET | /api/orders/:id | 查询单个订单 |
| POST | /api/orders | 创建订单（需验证用户） |
| PUT | /api/orders/:id/status | 更新订单状态 |
| GET | /api/users/:id/orders | 查询用户订单 |

## 快速开始

```bash
cd examples/microservice-architecture/order-service
go mod tidy
go run main.go
```

> 注意：创建订单前需确保用户服务已启动，因为会调用用户服务验证用户存在性。

## 数据模型

```json
{
  "id": 1,
  "user_id": 1,
  "product_ids": [1, 2],
  "total_amount": 99.9,
  "status": "pending",
  "created_at": "2026-06-19T10:00:00Z",
  "updated_at": "2026-06-19T10:00:00Z"
}
```

## 服务间通信

订单服务通过 HTTP 调用用户服务：

```
Order Service ──HTTP──► User Service (:8083)
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.New()` | 创建 IoC 容器 |
| `container.Register()` | 注册服务 Bean |
| `ginadapter.New()` | Gin 与容器集成 |
| HTTP Client | 服务间调用验证 |