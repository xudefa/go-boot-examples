# 任务管理系统示例

使用 go-boot 框架构建的任务管理系统，演示 IoC 容器、事件驱动、定时任务调度和 Gin Web 集成的综合应用。

## 功能特性

- **用户管理**：用户注册与查询
- **任务 CRUD**：创建、查询、更新任务
- **状态管理**：待办、进行中、已完成、已取消
- **优先级**：低、中、高、紧急
- **任务分配**：指定任务负责人
- **逾期检测**：自动识别逾期任务
- **定时提醒**：自动检查即将到期和逾期的任务
- **报告生成**：每日任务报告、用户任务报告
- **事件驱动**：任务创建、状态变更事件发布

## 架构设计

```
UserService ──► TaskService ◄── TaskReportService
                    │
                    ▼
            ReminderService ──► EventBus
                    │
                    ▼
            Scheduler (定时任务)
```

## 数据模型

| 模型 | 关键字段 |
|------|----------|
| User | ID, Username, Email, Role |
| Task | ID, Title, Description, AssigneeID, Status, Priority, DueDate |
| TaskReport | ID, TaskID, GeneratedAt, Content |

## API 端点

### 用户

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/users | 创建用户 |
| GET | /api/users/:id | 查询用户 |

### 任务

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/tasks | 查询所有任务 |
| GET | /api/tasks/:id | 查询单个任务 |
| POST | /api/tasks | 创建任务 |
| PUT | /api/tasks/:id | 更新任务 |
| PUT | /api/tasks/:id/status | 更新任务状态 |
| GET | /api/users/:id/tasks | 查询用户任务 |
| GET | /api/tasks/overdue | 查询逾期任务 |

### 报告

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/reports/daily | 生成每日报告 |
| GET | /api/reports/user/:id | 生成用户报告 |

## 快速开始

```bash
cd examples/task-management
go mod tidy
go run main.go
```

服务启动后访问 http://localhost:8082

## 使用示例

### 创建用户

```bash
curl -X POST http://localhost:8082/api/users \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","email":"alice@example.com"}'
```

### 创建任务

```bash
curl -X POST http://localhost:8082/api/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Fix bug","description":"Fix login issue","assignee_id":1,"creator_id":1,"priority":"high","due_date":"2026-06-20T18:00:00Z"}'
```

### 更新任务状态

```bash
curl -X PUT http://localhost:8082/api/tasks/1/status \
  -H 'Content-Type: application/json' \
  -d '{"status":"in_progress"}'
```

### 查询逾期任务

```bash
curl http://localhost:8082/api/tasks/overdue
```

## 定时任务

| 任务 | Cron 表达式 | 频率 | 说明 |
|------|-------------|------|------|
| check_due_tasks | */5 * * * * ? | 每 5 分钟 | 检查即将到期的任务 |
| check_overdue_tasks | 0 * * * * ? | 每小时 | 检查逾期任务 |
| generate_daily_report | 0 9 * * * ? | 每天 9:00 | 生成每日报告 |

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.New()` | 创建 IoC 容器 |
| `container.Register()` | 注册服务 Bean |
| `event.EventBus` | 事件发布/订阅 |
| `schedule.NewScheduler()` | 定时任务调度器 |
| `ginadapter.New()` | Gin 与容器集成 |