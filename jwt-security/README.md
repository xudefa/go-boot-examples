# JWT 安全集成示例

演示 go-boot 框架的安全模块和 JWT 模块的集成，实现基于 JWT 的认证和授权。

## 功能特性

- 同时启用安全模块和 JWT 模块
- 配置用户信息（内存存储）
- 登录获取 JWT Token
- 使用 Token 访问受保护的 API
- 基于角色的访问控制
- 配置排除路径（不需要认证的接口）

## 快速开始

```bash
cd examples/jwt-security
go mod tidy
go run .
```

## API 端点

| 端点 | 方法 | 需要认证 | 说明 |
|------|------|---------|------|
| `/login` | POST | 否 | 登录获取 Token |
| `/health` | GET | 否 | 健康检查 |
| `/api/public` | GET | 否 | 公开 API |
| `/api/protected` | GET | 是 | 受保护 API（需要 JWT） |
| `/api/admin` | GET | 是 | 管理员 API（需要 ADMIN 角色） |

## 使用示例

### 登录获取 Token

```bash
curl -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}'
```

### 访问受保护 API

```bash
curl http://localhost:8080/api/protected \
  -H 'Authorization: Bearer <token>'
```

### 访问管理员 API

```bash
curl http://localhost:8080/api/admin \
  -H 'Authorization: Bearer <token>'
```

## 安全规则

| 规则表达式 | 说明 |
|-----------|------|
| `permitAll` | 允许所有人访问 |
| `authenticated` | 仅允许已认证用户访问 |
| `hasRole('ADMIN')` | 仅允许具有 ADMIN 角色的用户访问 |
| `hasAnyRole('ADMIN','USER')` | 允许具有任一指定角色的用户访问 |

## 测试用户

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | ROLE_ADMIN, ROLE_USER |
| user | user123 | ROLE_USER |

## 核心组件

| 组件 | 说明 |
|------|------|
| `JwtUtil` | JWT 工具类，生成和验证 Token |
| `JwtAuthenticationFilter` | JWT 认证过滤器 |
| `SecurityFilterChain` | 安全过滤器链 |
| `AuthenticationManager` | 认证管理器 |
| `UserDetailsService` | 用户详情服务 |

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| JWT 认证 | Token 生成与验证 |
| 角色权限控制 | 基于角色的访问控制 |
| 安全过滤器链 | 统一安全处理 |
| 排除路径配置 | 公开接口配置 |