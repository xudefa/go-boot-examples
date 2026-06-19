# 电商系统示例

使用 go-boot 框架构建的完整电商系统，演示 IoC 容器、AOP 日志切面、事件驱动、定时任务和 Gin Web 集成的综合应用。

## 功能特性

- **用户管理**：用户注册、查询、角色管理（admin/customer）
- **商品管理**：商品 CRUD、分类、库存管理
- **购物车**：添加商品、查看购物车、移除商品
- **订单处理**：创建订单、支付、发货、状态跟踪
- **库存监控**：定时检查低库存商品并发布警报事件
- **AOP 日志**：使用 Around 通知记录商品操作日志
- **定时任务**：每分钟自动检查库存水位

## 架构设计

```
UserService ──┐
              │
ProductService ──► OrderService ──► CartService
      │              ▲
      │ (AOP)        │
      ▼              │
  ProxyFactory       │
                     │
EventBus ◄── InventoryMonitorService (定时任务)
```

## 数据模型

| 模型 | 关键字段 |
|------|----------|
| User | ID, Username, Email, Role |
| Product | ID, Name, Price, Stock, Category |
| Cart | UserID, Items[] |
| Order | ID, UserID, Items[], TotalAmount, Status |

## API 端点

### 用户

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/users | 创建用户 |
| GET | /api/users/:id | 查询用户 |

### 商品

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/products | 查询所有商品 |
| GET | /api/products/:id | 查询单个商品 |
| POST | /api/products | 创建商品 |

### 购物车

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/cart/:userId | 查看购物车 |
| POST | /api/cart/add | 添加到购物车 |
| DELETE | /api/cart/remove | 从购物车移除 |

### 订单

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/orders | 创建订单 |
| GET | /api/orders/:id | 查询订单 |
| GET | /api/users/:id/orders | 查询用户订单 |
| PUT | /api/orders/:id/status | 更新订单状态 |

### 库存

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/inventory/low-stock | 查询低库存商品 |

## 快速开始

```bash
cd examples/e-commerce-system
go mod tidy
go run main.go
```

服务启动后访问 http://localhost:8081

## 使用示例

### 创建用户

```bash
curl -X POST http://localhost:8081/api/users \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","email":"alice@example.com","role":"customer"}'
```

### 创建商品

```bash
curl -X POST http://localhost:8081/api/products \
  -H 'Content-Type: application/json' \
  -d '{"name":"Go Book","price":59.9,"stock":100,"category":"books"}'
```

### 创建订单

```bash
curl -X POST http://localhost:8081/api/orders \
  -H 'Content-Type: application/json' \
  -d '{"user_id":1,"items":[{"product_id":1,"quantity":2}]}'
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.New()` | 创建 IoC 容器 |
| `aop.NewProxyFactory()` | AOP 代理工厂 |
| `aop.Around()` | 环绕通知（日志切面） |
| `event.EventBus` | 事件发布/订阅 |
| `schedule.NewScheduler()` | 定时任务调度器 |
| `ginadapter.New()` | Gin 与容器集成 |