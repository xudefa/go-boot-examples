# 单例与原型作用域示例

演示 go-boot 核心中单例（Singleton）和原型（Prototype）两种 Bean 作用域的区别。

## 功能特性

- 注册单例 Bean（默认作用域，全局共享一个实例）
- 注册原型 Bean（每次获取都创建新实例）
- 验证单例 Bean 多次获取为同一实例
- 验证原型 Bean 每次获取为不同实例
- 修改单例实例影响所有引用，修改原型实例互相独立

## 快速开始

```bash
cd examples/core-scope
go mod tidy
go run .
```

## 预期输出

```
Container created

Registered singleton bean: singletonCounter
Singleton bean behavior:
  First get:  value = 0
  Second get: value = 100 (same instance, same value)
  Same instance: true

Registered prototype bean: prototypeCounter
Prototype bean behavior:
  First get:  value = 0
  Second get: value = 0
  Same instance: false (different instances)

After modifying first instance:
  First instance:  500
  Second instance: 0 (unchanged)

core-scope example completed successfully!
```

## 代码结构

```go
// 注册单例 Bean（默认）
container.Register("singletonCounter", core.Bean(&Counter{}))

// 注册原型 Bean
container.Register("prototypeCounter", core.Bean(&Counter{}), core.Prototype())
```

## 作用域对比

| 特性 | Singleton | Prototype |
|------|-----------|-----------|
| 实例数量 | 全局唯一 | 每次获取创建新实例 |
| 状态共享 | 所有引用共享状态 | 实例间状态独立 |
| 适用场景 | 无状态服务、配置 | 有状态对象、请求上下文 |

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `core.Bean()` | 注册单例 Bean（默认） |
| `core.Prototype()` | 注册原型 Bean |
| 作用域管理 | 控制 Bean 生命周期 |