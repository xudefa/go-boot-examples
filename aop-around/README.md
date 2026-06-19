# AOP Around 通知示例

演示 go-boot AOP 模块的 Around 通知——在方法执行前后进行环绕增强，可修改返回值。

## 功能特性

- 使用 `aop.Around()` 创建环绕通知
- 通过 `ProceedFunc` 控制目标方法是否执行
- 在方法执行前后添加横切逻辑（如性能计时）
- 修改目标方法的返回值
- 结合切点匹配器精确控制拦截范围

## 快速开始

```bash
cd examples/aop-around
go mod tidy
go run .
```

## 预期输出

```
=== AOP Around Advice Example ===

Created pointcut: MatchByName("Add")
Created Around advice (with timing and result modification)
Weaved target object

Calling Add(10, 20) on proxy:
  [Around] Before: Add, args=[10 20]
  [Target] Add(10, 20) called
  [Around] After: result=30, elapsed=50.xxxms
Final result: 60 (note: result was doubled by Around advice)

aop-around example completed successfully!
```

## 代码结构

```go
// 创建 Around 通知
advice := aop.Around(func(ctx *aop.Context, proceed aop.ProceedFunc) (any, error) {
    fmt.Printf("[Around] Before: %s, args=%v\n", ctx.Method.Name(), ctx.Args)
    
    // 执行目标方法
    result, err := proceed(ctx)
    
    fmt.Printf("[Around] After: result=%v, elapsed=%v\n", result, time.Since(start))
    
    // 修改返回值（例如：结果翻倍）
    if num, ok := result.(int); ok {
        return num * 2, err
    }
    return result, err
})

// 创建切面并织入
aspect := &aop.AspectMeta{
    Pointcut: aop.MatchByName("Add"),
    Advice:   advice,
}

weaver := aop.NewWeaver()
weaver.AddAspect(aspect)
proxy := weaver.Weave(target)
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `Around()` | 环绕通知 |
| `ProceedFunc` | 控制目标方法执行 |
| 返回值修改 | 拦截并修改返回值 |
| 性能计时 | 方法执行时间统计 |