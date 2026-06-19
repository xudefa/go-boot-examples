# 核心集成示例

演示 go-boot 核心模块的完整用法，包括依赖管理、条件注册、工厂创建和初始化回调。

## 功能特性

- 注册多个 Bean 并管理它们之间的依赖关系
- 使用 `core.DependsOn()` 显式声明 Bean 依赖
- 使用 `core.Condition()` 条件注册（依赖是否满足）
- 使用 `core.Init()` 注册初始化回调函数
- 使用 `core.Factory()` 工厂函数创建复杂 Bean
- 使用 `container.Invoke()` 方法注入已有的 Bean
- 字段注入自动解析

## 快速开始

```bash
cd examples/core-integrated
go mod tidy
go run .
```

## 预期输出

```
=== Core Integrated Example ===

Registered: logger
Registered: database (with condition)
Registered: userService (with dependencies and init)

Retrieved userService:
  Name: InitializedService
  Logger: true (injected)
  Database: true (injected, connected=true)

[INFO] AppLogger: User service is working!
[ERROR] AppLogger: This is an error message

Using container.Invoke():
[INFO] AppLogger: Invoked with db URL: localhost:5432

core-integrated example completed successfully!
```

## 代码结构

```go
// 条件注册
container.Register("database", core.Bean(db),
    core.Condition(func(c *core.Container) bool {
        return true // 条件满足时注册
    }),
)

// 依赖声明 + 初始化回调
container.Register("userService", core.Bean(&UserService{}),
    core.DependsOn("logger", "database"),
    core.Init(func(s *UserService) {
        s.Name = "InitializedService"
    }),
)

// 方法注入
container.Invoke(func(logger *Logger, db *Database) {
    logger.Info("Invoked with db URL: " + db.URL)
})
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.DependsOn()` | 显式声明 Bean 依赖 |
| `core.Condition()` | 条件注册 |
| `core.Init()` | 初始化回调 |
| `core.Factory()` | 工厂函数创建 Bean |
| `container.Invoke()` | 方法注入 |