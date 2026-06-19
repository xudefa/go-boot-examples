# GORM + Redis 双层缓存示例

演示 go-boot 框架中 GORM 数据库与 Redis 缓存的组合使用，实现 Cache-Aside 缓存模式。

## 功能特性

- GORM 数据库操作 + Redis 缓存
- Cache-Aside 模式（先读缓存 → 未命中读 DB → 回填缓存）
- Getter 回调模式
- 写操作时缓存失效（更新/删除）
- 使用 `boot.NewApplication()` 管理全生命周期

## 快速开始

```bash
cd examples/combine-db-cache
go mod tidy
go run main.go
```

## 关键概念

| 概念 | 说明 |
|------|------|
| Cache-Aside 模式 | 读：先查缓存，未命中则查 DB 并回填；写：更新 DB 后删除缓存 |
| `boot.NewApplication()` | 创建应用实例，管理全生命周期 |
| `environment.NewMapPropertySource()` | 基于 Map 的属性源配置 |
| `ggorm.NewRepository[Article]()` | 创建泛型 Repository |
| `redis.NewCacheWithConfig()` | 创建 Redis 缓存实例 |

## 缓存模式

```
读操作:
  1. 查询 Redis 缓存
  2. 缓存命中 → 返回数据
  3. 缓存未命中 → 查询 MySQL
  4. 将 DB 结果写入 Redis 缓存
  5. 返回数据

写操作:
  1. 更新 MySQL
  2. 删除 Redis 缓存（保证下次读取获取最新数据）
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| GORM 集成 | 数据库 CRUD 操作 |
| Redis 缓存 | 高性能缓存层 |
| Cache-Aside | 经典缓存模式实现 |
| `boot.NewApplication()` | 应用生命周期管理 |