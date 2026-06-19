# JWT 工具示例

演示 go-boot 框架的 JWT 工具，包括 Token 生成、解析、验证和刷新等功能。

## 功能特性

- 创建 JWT 工具实例并自定义配置
- 生成 Access Token 和 Refresh Token
- 解析 Token 获取声明信息
- 验证 Token 有效性
- 刷新过期 Token
- 获取 Token 中的用户信息
- 获取 Token 剩余有效期

## 快速开始

```bash
cd examples/jwt-example
go mod tidy
go run .
```

## 关键概念

| 概念 | 说明 |
|------|------|
| `jwt.NewJWTUtil()` | 创建 JWT 工具实例，支持函数式选项配置 |
| `GenerateToken()` | 生成 Access Token 和 Refresh Token |
| `ParseToken()` | 解析 Token 获取完整的声明信息 |
| `ValidateToken()` | 验证 Token 是否有效（包括过期检查） |
| `RefreshToken()` | 使用现有 Token 生成新的 Token 对 |
| `GetSubject()` | 获取 Token 中的用户名（Subject） |
| `GetUserId()` | 获取 Token 中的用户 ID |
| `GetClaims()` | 获取 Token 的完整声明信息 |
| `GetRemainingTime()` | 获取 Token 的剩余有效期 |

## 配置选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `WithSecretKey()` | 设置 JWT 签名密钥 | `go-bootJwtSecret` |
| `WithIssuer()` | 设置签发者 | `go-boot` |
| `WithExpiresDuration()` | 设置 Access Token 过期时间 | `10 分钟` |
| `WithRefreshExpiresDuration()` | 设置 Refresh Token 过期时间 | `1 小时` |

## 使用示例

### 创建 JWT 工具

```go
jwtUtil := jwt.NewJWTUtil(
    jwt.WithSecretKey("my-secret-key"),
    jwt.WithIssuer("my-app"),
    jwt.WithExpiresDuration(30*time.Minute),
    jwt.WithRefreshExpiresDuration(2*time.Hour),
)
```

### 生成 Token

```go
tokens, err := jwtUtil.GenerateToken("alice", "user-id-123")
fmt.Println("Access Token:", tokens.AccessToken)
fmt.Println("Refresh Token:", tokens.RefreshToken)
```

### 验证 Token

```go
valid := jwtUtil.ValidateToken(tokens.AccessToken)
```

### 刷新 Token

```go
newTokens, err := jwtUtil.RefreshToken(tokens.RefreshToken)
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| JWT 工具 | Token 生成与验证 |
| 函数式选项 | 灵活的配置方式 |
| Token 刷新 | 无感续期机制 |
| 声明解析 | 获取用户信息 |