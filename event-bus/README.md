# 事件总线示例

演示 go-boot 框架的事件驱动支持，包括事件发布/订阅、自定义事件类型和内置生命周期事件。

## 功能特性

- **事件总线**：`NewEventBus()` 创建发布/订阅中心
- **事件订阅**：通过 `Subscribe()` 注册事件处理器
- **事件发布**：通过 `Publish()` 广播事件到所有订阅者
- **自定义事件**：自定义事件类型（如 `UserRegisteredEvent`）
- **内置事件**：框架生命周期事件

## 内置事件

| 事件类型 | 触发时机 |
|----------|----------|
| `EventEnvironmentPrepared` | 环境准备完成 |
| `EventContextRefreshed` | 上下文刷新 |
| `EventApplicationStarted` | 应用开始启动 |
| `EventApplicationReady` | 应用就绪 |
| `EventApplicationStopped` | 应用停止 |

## 快速开始

```bash
cd examples/event-bus
go mod tidy
go run main.go
```

## 使用示例

### 创建事件总线

```go
eventBus := event.NewEventBus()
```

### 订阅事件

```go
eventBus.Subscribe("user.registered", func(e event.ApplicationEvent) {
    fmt.Println("User registered:", e.Type())
})
```

### 发布事件

```go
eventBus.Publish(&UserRegisteredEvent{
    EventType: "user.registered",
    EventTime: time.Now(),
    UserID:    1,
})
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `event.NewEventBus()` | 创建事件总线 |
| `Subscribe()` | 注册事件处理器 |
| `Publish()` | 发布事件 |
| `ApplicationEvent` 接口 | 自定义事件类型 |