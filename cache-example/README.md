# 缓存抽象与内存缓存示例

演示 go-boot 框架的缓存抽象层，包括内存缓存和多种缓存操作。

## 功能特性

- 创建内存缓存实例
- 使用 `Get`、`Set`、`Exists`、`TTL` 等基础操作
- 使用批量操作 `GetMulti` / `SetMulti`
- 使用 `GetWithGetter` 实现延迟加载/缓存穿透保护

## 快速开始

```bash
cd examples/cache-example
go mod tidy
go run main.go
```

## 关键概念

| 概念 | 说明 |
|------|------|
| `cache.NewMemoryCache()` | 创建基于内存的缓存实例 |
| `cache.Cache` 接口 | 统一缓存抽象，定义 `Get`、`Set`、`Exists`、`Del`、`Clear` 等方法 |
| `GetWithGetter` | 缓存未命中时通过回调函数自动加载数据（缓存穿透保护） |
| `GetMulti` / `SetMulti` | 批量操作，减少网络开销 |
| `cache.Getter` | 数据加载回调函数类型 |

## 使用示例

### 基础操作

```go
cache := cache.NewMemoryCache()
cache.Set("key", "value", 5*time.Minute)
val, _ := cache.Get("key")
```

### 批量操作

```go
cache.SetMulti(map[string]any{
    "key1": "value1",
    "key2": "value2",
}, 5*time.Minute)
```

### 延迟加载

```go
val, _ := cache.GetWithGetter("key", func() (any, error) {
    return loadDataFromDB()
})
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| 缓存抽象层 | 统一的 `cache.Cache` 接口 |
| 内存缓存 | 基于 map 的高性能本地缓存 |
| 缓存穿透保护 | `GetWithGetter` 回调机制 |
| 批量操作 | 减少多次调用的开销 |