# Gin + GORM 集成

演示 go-boot 与 Gin + GORM 的完整集成，实现 RESTful CRUD API。

## 功能特性

- 使用 `boot.NewApplication()` 创建应用实例
- 连接 MySQL 数据库（无数据库时自动进入演示模式）
- 自动迁移数据表
- 泛型 `Repository[T]` 实现类型安全 CRUD
- RESTful API：用户增删改查

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/users | 查询所有用户 |
| GET | /api/users/:id | 查询单个用户 |
| POST | /api/users | 创建新用户 |
| DELETE | /api/users/:id | 删除用户 |

## 快速开始

```bash
cd examples/gin-integrated
go mod tidy
go run main.go
```

无数据库配置时，自动进入演示模式并打印路由结构。

## 使用示例

### 创建用户

```bash
curl -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

### 查询所有用户

```bash
curl http://localhost:8080/api/users
```

### 删除用户

```bash
curl -X DELETE http://localhost:8080/api/users/1
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `boot.NewApplication()` | 创建应用实例 |
| `ginadapter.New()` | Gin 与容器集成 |
| `ggorm.NewRepository[T]()` | 泛型 Repository |
| `gorm.AutoMigrate` | 自动建表 |
| RESTful 设计 | 标准 HTTP 方法映射 |

## 依赖

- [gin-gonic/gin](https://github.com/gin-gonic/gin) — HTTP 框架
- [gorm.io/gorm](https://gorm.io) — ORM 框架
- `github.com/xudefa/go-boot-gin/server` — Gin 集成
- `github.com/xudefa/go-boot-gorm` — GORM 集成
- `github.com/xudefa/go-boot/boot` — 应用启动
- `github.com/xudefa/go-boot/core` — IoC 容器