# Redis 缓存示例

演示 go-boot 的 Redis 集成模块，实现 Redis 缓存功能及不可用时的内存缓存回退。

## 功能特性

- 创建 Redis 缓存实例，支持函数式选项配置
- Redis 不可用时自动回退到内存缓存
- 使用统一的 `cache.Cache` 接口进行缓存操作
- 支持自定义缓存前缀和默认 TTL

## 快速开始

```bash
cd examples/redis-cache
go mod tidy
go run main.go
```

## 关键概念

| 概念 | 说明 |
|------|------|
| `redis.NewCacheWithConfig()` | 创建 Redis 缓存实例，支持函数式选项 |
| `redis.WithPrefix()` | 设置缓存键前缀 |
| `redis.WithDefaultTTL()` | 设置默认过期时间 |
| `cache.Cache` 接口 | 统一的缓存抽象，支持 `Get`、`Set`、`Exists`、`Del`、`Clear` |
| `cache.NewMemoryCache()` | 创建内存缓存（Redis 不可用时的回退方案） |

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| Redis 集成 | 高性能分布式缓存 |
| 自动回退 | Redis 不可用时降级到内存缓存 |
| 统一接口 | `cache.Cache` 抽象层 |
| 函数式选项 | 灵活的配置方式 |