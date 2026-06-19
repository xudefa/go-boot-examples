# API 网关

微服务架构的统一入口，负责将请求路由到对应的后端服务。

## 端口

`http://localhost:8080`

## 路由规则

| 路径前缀 | 转发目标 | 说明 |
|----------|----------|------|
| /users/* | User Service (:8083) | 用户服务路由 |
| /orders/* | Order Service (:8084) | 订单服务路由 |
| /health | 网关自身 | 网关健康检查 |
| / | 网关自身 | 服务信息页 |

## 快速开始

```bash
cd examples/microservice-architecture/api-gateway
go mod tidy
go run main.go
```

> 注意：网关启动前需确保后端服务（User Service、Order Service）已启动。

## 使用示例

### 查看服务信息

```bash
curl http://localhost:8080/
```

### 健康检查

```bash
curl http://localhost:8080/health
```

### 转发到用户服务

```bash
curl http://localhost:8080/users/api/users
```

### 转发到订单服务

```bash
curl http://localhost:8080/orders/api/orders
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.New()` | 创建 IoC 容器 |
| `container.Register()` | 注册服务 Bean |
| `ginadapter.New()` | Gin 与容器集成 |
| HTTP 代理 | 请求转发与响应回传 |
| 路由匹配 | 路径前缀匹配转发 |