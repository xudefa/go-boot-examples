# GORM 基础使用示例

演示 go-boot 的 GORM 集成模块，包括数据库连接配置、泛型 Repository 和 CRUD 操作。

## 功能特性

- 连接 MySQL 数据库（支持 PostgreSQL、SQLServer、SQLite）
- 使用函数式选项配置连接参数
- 泛型 `Repository[T]` 提供类型安全的 CRUD
- 启动器模式（Starter）实现连通性检查和自动迁移

## 快速开始

```bash
cd examples/gorm-basic
go mod tidy
go run main.go
```

## 关键概念

| 概念 | 说明 |
|------|------|
| `gorm.OpenMySQL()` | 使用函数式选项创建 MySQL 连接 |
| `gorm.WithHost()` / `gorm.WithPort()` | 配置数据库连接参数 |
| `gorm.Config.DSNForMySQL()` | 自动生成 DSN 连接串 |
| `Repository[T]` | 泛型数据仓库，类型安全 CRUD |
| `gorm.NewStarter()` | 启动器，包含连通性检查 |
| `gorm.NewAutoMigrateStarter()` | 启动器，包含自动迁移 |

## 支持的数据库

| 数据库 | 方法 |
|--------|------|
| MySQL | `gorm.OpenMySQL()` |
| PostgreSQL | `gorm.OpenPostgreSQL()` |
| SQLServer | `gorm.OpenSQLServer()` |
| SQLite | `gorm.OpenSQLite()` |

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| 函数式选项模式 | 灵活的连接配置 |
| 泛型 Repository | 类型安全的数据库操作 |
| Starter 模式 | 生命周期管理与连通性检查 |