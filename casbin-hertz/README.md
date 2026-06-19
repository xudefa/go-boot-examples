# Casbin + Hertz 访问控制示例

演示 Casbin 权限管理器与字节跳动 Hertz HTTP 框架的集成，实现基于角色的访问控制（RBAC）。

## 功能特性

- **RBAC 权限模型**：基于角色的访问控制
- **Casbin 授权中间件**：在 Hertz 路由上自动进行权限检查
- **请求头提取用户**：从 `X-User` 请求头获取当前用户
- **路径与方法映射**：自动将请求路径和方法映射为操作类型（read/write）

## 快速开始

```bash
cd examples/casbin-hertz
go mod tidy
go run .
```

服务启动在 `:8091` 端口。

## API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/data1 | 读取 data1 |
| POST | /api/data1/write | 写入 data1 |
| GET | /api/data2 | 读取 data2 |
| POST | /api/data2/write | 写入 data2 |

## 使用示例

```bash
# Alice 读取 data1 — 应返回 200
curl -H 'X-User: alice' http://localhost:8091/api/data1

# Bob 写入 data1 — 应返回 403
curl -X POST -H 'X-User: bob' http://localhost:8091/api/data1/write
```

## 策略配置

策略文件位于 `policy.csv`，使用 RBAC 模型：

| 用户 | 角色 | 资源 | 操作 |
|------|------|------|------|
| alice | admin | data1, data2 | read, write |
| bob | user | data1, data2 | read |

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| Casbin 集成 | 权限管理框架 |
| RBAC 模型 | 基于角色的访问控制 |
| Hertz 中间件 | HTTP 请求拦截 |
| 请求头提取 | 从请求头获取用户信息 |

## 依赖

- `github.com/xudefa/go-boot/casbin` — 权限管理
- `github.com/xudefa/go-boot/hertz` — Hertz 服务器集成
- `github.com/cloudwego/hertz` — 字节跳动 HTTP 框架