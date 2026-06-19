# AOP 多通知组合与排序示例

演示 go-boot AOP 模块中多种通知类型的组合使用，以及通过 Order 控制执行顺序。

## 功能特性

- 使用 Before 通知（前置处理）
- 使用 Around 通知（环绕增强，支持性能计时）
- 使用 After 通知（后置处理，无论是否异常都会执行）
- 使用 AfterReturning 通知（方法成功返回后执行）
- 通过 `AspectMeta.Order` 字段控制通知执行顺序
- 多个通知形成链式调用

## 快速开始

```bash
cd examples/aop-integrated
go mod tidy
go run .
```

## 预期输出

```
=== AOP Integrated Example (Multiple Advices) ===

Created proxy with 4 advices (orders 1-4)

Calling CreateOrder(100, 50.0):
[Before-1] Method: CreateOrder, Args: [100 50]
[Around-2] Before proceed: CreateOrder
  [Target] CreateOrder called
[Around-2] After proceed: result=Order-100-50, elapsed=30.xxxms
[After-3] Method completed: CreateOrder
[AfterReturning-4] Method returned: Order-100-50

Final result: Order-100-50

aop-integrated example completed successfully!
```

## 通知执行顺序

| Order | 通知类型 | 执行时机 |
|-------|----------|----------|
| 1 | Before | 方法执行前 |
| 2 | Around | 方法执行前后 |
| 3 | After | 方法执行后（无论是否异常） |
| 4 | AfterReturning | 方法成功返回后 |

## 代码结构

```go
// 创建多个通知
advices := []aop.Advice{
    aop.Before(func(ctx *aop.Context) {
        fmt.Printf("[Before-1] Method: %s\n", ctx.Method.Name())
    }),
    aop.Around(func(ctx *aop.Context, proceed aop.ProceedFunc) (any, error) {
        fmt.Printf("[Around-2] Before proceed: %s\n", ctx.Method.Name())
        result, err := proceed(ctx)
        fmt.Printf("[Around-2] After proceed: result=%v\n", result)
        return result, err
    }),
    aop.After(func(ctx *aop.Context) {
        fmt.Printf("[After-3] Method completed: %s\n", ctx.Method.Name())
    }),
    aop.AfterReturning(func(ctx *aop.Context, result any) {
        fmt.Printf("[AfterReturning-4] Method returned: %v\n", result)
    }),
}

// 设置 Order
for i, advice := range advices {
    advice.Meta().Order = i + 1
}
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| Before 通知 | 方法执行前拦截 |
| Around 通知 | 方法执行前后环绕增强 |
| After 通知 | 方法执行后处理 |
| AfterReturning 通知 | 方法成功返回后处理 |
| Order 排序 | 控制通知执行顺序 |