# 工厂函数示例

演示 go-boot 核心的工厂函数机制，使用 `core.Factory()` 动态创建 Bean 实例。

## 功能特性

- 使用 `core.Factory()` 注册工厂函数
- 工厂函数接收容器引用，可访问其他 Bean
- 结合 `reflect.TypeFor[]` 注册工厂 Bean 的类型
- 支持与 `core.Singleton()` 等作用域选项组合使用
- 工厂函数可执行初始化逻辑后再返回实例

## 快速开始

```bash
cd examples/core-factory
go mod tidy
go run .
```

## 预期输出

```
Container created

Registering factory-created bean...
  Factory: creating Config...
Retrieved Config: Env=production, Port=8080

core-factory example completed successfully!
```

## 代码结构

```go
// 注册工厂函数
container.Register("config", core.Factory(func(c *core.Container) *Config {
    fmt.Println("Factory: creating Config...")
    return &Config{
        Env:  "production",
        Port: 8080,
    }
}, reflect.TypeFor[*Config]()))
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.Factory()` | 注册工厂函数 |
| 容器引用 | 工厂函数可访问其他 Bean |
| 类型注册 | `reflect.TypeFor[]` 指定类型 |
| 初始化逻辑 | 工厂函数执行初始化 |