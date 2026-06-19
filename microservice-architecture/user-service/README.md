# 用户服务

微服务架构中的用户管理服务，负责用户的创建、查询和更新。

## 端口

`http://localhost:8083`

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /health | 健康检查 |
| GET | /api/users | 查询所有用户 |
| GET | /api/users/:id | 查询单个用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/:id | 更新用户 |

## 快速开始

```bash
cd examples/microservice-architecture/user-service
go mod tidy
go run main.go
```

## 数据模型

```json
{
  "id": 1,
  "username": "alice",
  "email": "alice@example.com",
  "phone": "1234567890",
  "address": "Beijing",
  "created": "2023-01-01"
}
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.New()` | 创建 IoC 容器 |
| `container.Register()` | 注册服务 Bean |
| `ginadapter.New()` | Gin 与容器集成 |