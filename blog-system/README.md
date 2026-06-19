# 博客系统示例

使用 go-boot 框架构建的完整博客系统，演示 IoC 容器、依赖注入、事件驱动和 Gin Web 集成的综合应用。

## 功能特性

- **用户管理**：用户注册与查询
- **文章管理**：创建、编辑、删除、查询文章
- **评论系统**：为文章添加评论
- **事件驱动**：文章创建时发布事件通知
- **IoC 容器**：所有服务注册到容器统一管理
- **依赖注入**：服务间通过构造函数注入依赖
- **HTML 模板**：使用 Gin 加载 HTML 模板渲染首页

## 架构设计

```
UserService ──┐
              ├──► ArticleService ──┐
              │                     │
              └──► CommentService ◄─┘
                                    │
EventBus ◄──────── BlogEventService┘
```

## 数据模型

| 模型 | 字段 |
|------|------|
| User | ID, Username, Email |
| Article | ID, Title, Content, Author, Created |
| Comment | ID, ArticleID, Content, Author, Created |

## API 端点

### 用户

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/users | 创建用户 |
| GET | /api/users/:id | 查询用户 |

### 文章

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/articles | 查询所有文章 |
| GET | /api/articles/:id | 查询单篇文章 |
| POST | /api/articles | 创建文章 |
| PUT | /api/articles/:id | 更新文章 |
| DELETE | /api/articles/:id | 删除文章 |

### 评论

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/articles/:id/comments | 查询文章评论 |
| POST | /api/comments | 创建评论 |

### 页面

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | / | 博客首页（HTML 模板） |

## 快速开始

```bash
cd examples/blog-system
go mod tidy
go run main.go
```

服务启动后访问 http://localhost:8080

## 使用示例

### 创建用户

```bash
curl -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","email":"alice@example.com"}'
```

### 创建文章

```bash
curl -X POST http://localhost:8080/api/articles \
  -H 'Content-Type: application/json' \
  -d '{"title":"Hello World","content":"My first post","author_id":1}'
```

### 创建评论

```bash
curl -X POST http://localhost:8080/api/comments \
  -H 'Content-Type: application/json' \
  -d '{"article_id":1,"content":"Great post!","author_id":1}'
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.New()` | 创建 IoC 容器 |
| `container.Register()` | 注册服务 Bean |
| `event.EventBus` | 事件发布/订阅 |
| `ginadapter.New()` | Gin 与容器集成 |
| 构造函数注入 | 服务间依赖传递 |