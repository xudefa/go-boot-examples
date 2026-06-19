# AOP 基础 Before 通知示例

演示 go-boot AOP 模块的基础用法——使用 Before 通知在目标方法执行前进行拦截。

## 功能特性

- 使用 `aop.MatchByName()` 创建切点匹配器
- 使用 `aop.Before()` 创建前置通知
- 创建 `AspectMeta` 组装切面和切点
- 使用 `aop.NewWeaver()` 创建织入器
- 使用 `weaver.Weave()` 织入切面到目标对象
- 在代理对象上调用方法触发通知

## 快速开始

```bash
cd examples/aop-basic
go mod tidy
go run .
```

## 预期输出

```
=== AOP Basic Example (Before Advice) ===

Created pointcut: MatchByName("GetUser")
Created Before advice
Created AspectMeta

Created weaver and added aspect

Weaved target object, got proxy

Calling method on proxy:
  [Before] Method: GetUser
  [Before] Args: [123]
  [Target] GetUser(123) called
Result: User123

aop-basic example completed successfully!
```

## 代码结构

```go
// 创建切点匹配器
pointcut := aop.MatchByName("GetUser")

// 创建 Before 通知
advice := aop.Before(func(ctx *aop.Context) {
    fmt.Printf("[Before] Method: %s\n", ctx.Method.Name())
    fmt.Printf("[Before] Args: %v\n", ctx.Args)
})

// 创建切面
aspect := &aop.AspectMeta{
    Pointcut: pointcut,
    Advice:   advice,
}

// 创建织入器并织入切面
weaver := aop.NewWeaver()
weaver.AddAspect(aspect)

// 织入目标对象，获取代理
proxy := weaver.Weave(target)

// 调用代理对象方法
result := proxy.GetUser(123)
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `MatchByName()` | 切点匹配器 |
| `Before()` | 前置通知 |
| `AspectMeta` | 切面元数据 |
| `NewWeaver()` | 创建织入器 |
| `Weave()` | 织入切面生成代理 |