# Gin Hello World

演示 go-boot 与 Gin 框架的基础集成，快速创建一个 HTTP 服务器并注册路由。

## 功能特性

- 使用 `ginadapter.New()` 创建 Gin 服务器
- 绑定到 `:8080` 端口
- 注册 `GET /hello` 路由返回 JSON 响应
- 将服务器实例注册到 IoC 容器

## 快速开始

```bash
cd examples/gin-hello
go mod tidy
go run main.go
```

启动后访问 http://localhost:8080/hello

## 预期响应

```json
{
  "message": "Hello from go-boot!"
}
```

## 代码结构

```go
// 创建 Gin 服务器并绑定到容器
g := ginadapter.New(ginadapter.WithContainer(container))

// 注册路由
g.GET("/hello", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Hello from go-boot!"})
})

// 启动服务器
g.Start()
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.New()` | 创建 IoC 容器 |
| `ginadapter.New()` | 创建 Gin 服务器并集成容器 |
| `container.Register()` | 注册服务器 Bean |
| 路由注册 | Gin 路由定义 |

## 依赖

- [gin-gonic/gin](https://github.com/gin-gonic/gin) — HTTP 框架
- `github.com/xudefa/go-boot-gin/server` — go-boot Gin 集成模块
- `github.com/xudefa/go-boot/core` — IoC 容器