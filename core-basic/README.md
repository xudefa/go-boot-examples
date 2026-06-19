# IoC 容器基础示例

演示 go-boot 核心模块的基础 Bean 注册、获取和容器操作。这是学习 go-boot 框架的起点。

## 功能特性

- 创建 IoC 容器实例
- 使用 `core.Bean()` 注册服务实例
- 使用 `container.Get()` 从容器中获取 Bean
- 使用 `container.Has()` 判断 Bean 是否存在
- 类型断言将 Bean 转换为具体类型

## 快速开始

```bash
cd examples/core-basic
go mod tidy
go run .
```

## 预期输出

```
Container created
Bean registered: userService
Retrieved bean: &{Name:MyService}
Container has userService: true

core-basic example completed successfully!
```

## 代码结构

```go
// 创建容器
container := core.New()

// 注册 Bean
container.Register("userService", core.Bean(&UserService{Name: "MyService"}))

// 获取 Bean
bean, _ := container.Get("userService")
service := bean.(*UserService)

// 检查 Bean 是否存在
exists := container.Has("userService")
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.New()` | 创建 IoC 容器 |
| `core.Bean()` | 注册 Bean 实例 |
| `container.Get()` | 获取 Bean |
| `container.Has()` | 检查 Bean 是否存在 |