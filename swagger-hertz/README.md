# Hertz + Swagger 集成示例

演示 go-boot 与 Hertz 框架及 Swagger API 文档的集成，快速创建带 API 文档的 HTTP 服务。

## 功能特性

- 与 go-boot IoC 容器集成
- Hertz Web 框架（字节跳动开源）
- Swagger UI 文档界面
- API 文档定义（内嵌 JSON）
- 请求/响应模型定义
- 安全认证配置

## 快速开始

```bash
cd examples/swagger-hertz
go mod tidy
go run .
```

启动后访问 http://localhost:8081/swagger/index.html

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /swagger/index.html | Swagger UI 界面 |
| GET | /swagger/doc.json | Swagger JSON 文档 |
| GET | /api/v1/hello | 获取欢迎信息 |
| GET | /api/v1/users/:id | 获取用户信息 |
| POST | /api/v1/users | 创建新用户 |

## 使用示例

### 测试 API

```bash
# 获取欢迎信息
curl http://localhost:8081/api/v1/hello

# 获取用户信息
curl http://localhost:8081/api/v1/users/1

# 创建新用户
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

## 代码结构

```go
// 创建 IoC 容器
container := core.New()

// 创建 Hertz 服务器
server := hs.NewServer(hs.WithHost(":8081"))

// 注册路由
api := server.Group("/api/v1")
{
    api.GET("/hello", helloHandler)
    api.GET("/users/:id", getUserHandler)
    api.POST("/users", createUserHandler)
}

// 注册到容器
container.Register("hertzServer", core.Bean(server))

// 启动服务
server.Start()
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.New()` | 创建 IoC 容器 |
| Hertz 服务器 | 高性能 HTTP 框架 |
| Swagger 集成 | API 文档自动生成 |
| 路由分组 | API 版本管理 |

## 依赖

- [cloudwego/hertz](https://github.com/cloudwego/hertz) — Hertz Web 框架
- `github.com/xudefa/go-boot/core` — IoC 容器