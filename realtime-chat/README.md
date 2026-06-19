# 实时聊天系统示例

使用 go-boot 框架构建的实时聊天系统，演示 IoC 容器、事件驱动、Gin Web 集成和 HTML 模板渲染。

## 功能特性

- **用户管理**：用户注册、在线状态跟踪
- **消息系统**：发送消息、消息历史记录
- **聊天室**：创建房间、加入房间
- **在线用户**：查看当前在线用户列表
- **系统消息**：用户上下线自动广播
- **事件驱动**：用户加入/离开事件发布
- **HTML 界面**：提供聊天 Web 界面

## 架构设计

```
ChatServer (核心)
  ├── UserService
  ├── MessageService
  └── RoomService
        │
        ▼
     EventBus (用户加入/离开事件)
```

## 数据模型

| 模型 | 关键字段 |
|------|----------|
| User | ID, Username, Online, JoinedAt |
| Message | ID, SenderID, Sender, Room, Content, Timestamp |
| Room | ID, Name, Description, Members |

## API 端点

### 用户

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/users | 注册用户 |
| GET | /api/users/:id | 查询用户 |
| GET | /api/users/online | 在线用户列表 |

### 消息

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/messages | 发送消息 |
| GET | /api/messages/recent | 获取最近消息 |

### 房间

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/rooms | 创建房间 |
| POST | /api/rooms/:id/join | 加入房间 |

### 页面

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | / | 聊天首页（HTML 模板） |

## 快速开始

```bash
cd examples/realtime-chat
go mod tidy
go run main.go
```

服务启动后访问 http://localhost:8085

## 使用示例

### 注册用户

```bash
curl -X POST http://localhost:8085/api/users \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice"}'
```

### 发送消息

```bash
curl -X POST http://localhost:8085/api/messages \
  -H 'Content-Type: application/json' \
  -d '{"sender_id":1,"room":"general","content":"Hello!"}'
```

### 获取最近消息

```bash
curl http://localhost:8085/api/messages/recent?count=20
```

### 创建房间

```bash
curl -X POST http://localhost:8085/api/rooms \
  -H 'Content-Type: application/json' \
  -d '{"name":"tech","description":"技术讨论"}'
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.New()` | 创建 IoC 容器 |
| `container.Register()` | 注册服务 Bean |
| `event.EventBus` | 事件发布/订阅 |
| `ginadapter.New()` | Gin 与容器集成 |
| Channel 通信 | 消息广播机制 |
| `sync.RWMutex` | 并发安全保护 |