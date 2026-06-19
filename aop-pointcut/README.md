# AOP 多种切点匹配器示例

演示 go-boot AOP 模块支持的多种切点匹配方式，精确控制通知的应用范围。

## 功能特性

- `MatchByName()` 精确匹配方法名称
- `MatchByNamePrefix()` 按方法名前缀匹配
- `MatchByRegex()` 使用正则表达式匹配方法名
- 不同切点可应用不同通知逻辑
- 每个切面独立织入，互不干扰

## 快速开始

```bash
cd examples/aop-pointcut
go mod tidy
go run .
```

## 预期输出

```
=== AOP PointCut Matchers Example ===

1. MatchByName("DoWork"):
   [Before] DoWork
  [Target] DoWork() called
  [Target] DoSomething(test) called

2. MatchByNamePrefix("Do"):
  [Target] DoWork() called
   [Before-Do] DoWork
  [Target] DoSomething(test) called
   [Before-Do] DoSomething
  [Target] HandleRequest(1) called

3. MatchByRegex("^Handle.*"):
  [Target] DoWork() called
  [Target] HandleRequest(1) called
   [Before-Handle] HandleRequest
  [Target] ProcessData(test) called

aop-pointcut example completed successfully!
```

## 切点匹配器对比

| 匹配器 | 说明 | 示例 |
|--------|------|------|
| `MatchByName()` | 精确匹配方法名 | `MatchByName("DoWork")` |
| `MatchByNamePrefix()` | 按方法名前缀匹配 | `MatchByNamePrefix("Do")` |
| `MatchByRegex()` | 正则表达式匹配 | `MatchByRegex("^Handle.*")` |

## 代码结构

```go
// 精确匹配
pointcut1 := aop.MatchByName("DoWork")

// 前缀匹配
pointcut2 := aop.MatchByNamePrefix("Do")

// 正则匹配
pointcut3 := aop.MatchByRegex("^Handle.*")

// 创建切面
aspect1 := &aop.AspectMeta{
    Pointcut: pointcut1,
    Advice:   aop.Before(func(ctx *aop.Context) {
        fmt.Printf("[Before] %s\n", ctx.Method.Name())
    }),
}
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| 精确匹配 | 指定方法名拦截 |
| 前缀匹配 | 批量拦截同名前缀方法 |
| 正则匹配 | 灵活的模式匹配 |
| 多切面织入 | 多个切面独立工作 |