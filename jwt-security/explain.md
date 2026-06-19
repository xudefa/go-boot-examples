=== JWT Security Integration Example ===

Application started successfully!
Server running on http://localhost:8080

Endpoints:
  POST   /login          - 登录获取 Token（不需要认证）
  GET    /health         - 健康检查（不需要认证）
  GET    /api/public     - 公开 API（不需要认证）
  GET    /api/protected  - 受保护 API（需要 JWT Token）
  GET    /api/admin      - 管理员 API（需要 JWT Token 和 ADMIN 角色）

Usage:
  1. 登录获取 Token:
     curl -X POST http://localhost:8080/login -H 'Content-Type: application/json' -d '{"username":"admin","password":"admin123"}'

  2. 使用 Token 访问受保护 API:
     curl http://localhost:8080/api/protected -H 'Authorization: Bearer <token>'
