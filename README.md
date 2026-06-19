# go-boot 示例项目

本目录包含了多个使用 go-boot 框架构建的示例项目，涵盖了从基础功能到完整业务系统的不同场景。

## 快速开始

### 运行示例

```bash
# 进入示例目录
cd examples/gin-hello

# 安装依赖并运行
go mod tidy
go run main.go
```

### 先决条件

- Go 1.21+
- 部分示例需要外部服务支持（Redis、MySQL、Nacos 等），详见各示例的 README

---

## 示例项目列表

### 核心模块 (Core)

| 示例 | 说明 | 关键特性 |
|------|------|----------|
| [core-basic](core-basic/) | 基础 Bean 注册与获取 | Bean 注册、依赖获取 |
| [core-injection](core-injection/) | 字段注入与依赖注入 | `inject` 标签、自动装配 |
| [core-factory](core-factory/) | 工厂函数创建 Bean | Factory 模式、延迟初始化 |
| [core-scope](core-scope/) | Bean 作用域管理 | Singleton、Prototype 作用域 |
| [core-integrated](core-integrated/) | 核心功能集成示例 | IoC + AOP + 生命周期综合演示 |
| [conditional](conditional/) | 条件注册示例 | 根据配置/环境动态注册 Bean |
| [property-source](property-source/) | 属性源示例 | 多来源配置、优先级管理 |

### AOP 模块

| 示例 | 说明 | 关键特性 |
|------|------|----------|
| [aop-basic](aop-basic/) | Before 通知基础示例 | 前置通知、方法拦截 |
| [aop-around](aop-around/) | Around 通知与 proceed 控制 | 环绕通知、执行链控制 |
| [aop-pointcut](aop-pointcut/) | 多种切点匹配器 | ByName、ByPrefix、ByRegex |
| [aop-integrated](aop-integrated/) | 多通知与排序示例 | 多切面协同、通知顺序 |

### Web 模块

| 示例 | 说明 | 关键特性 |
|------|------|----------|
| [gin-hello](gin-hello/) | 基础 Gin HTTP 服务器 | 路由注册、中间件、容器集成 |
| [gin-integrated](gin-integrated/) | Gin + GORM 集成示例 | Web + 数据库完整流程 |
| [swagger-hertz](swagger-hertz/) | Hertz + Swagger 文档集成 | API 文档自动生成 |

### 数据模块

| 示例 | 说明 | 关键特性 |
|------|------|----------|
| [gorm-basic](gorm-basic/) | GORM 基础操作与 Repository 模式 | CRUD、泛型 Repository |
| [combine-db-cache](combine-db-cache/) | 数据库 + 缓存集成示例 | Cache-Aside 模式、读写分离 |

### 配置模块

| 示例 | 说明 | 关键特性 |
|------|------|----------|
| [config-example](config-example/) | 基础配置加载 | 环境变量、配置文件 |
| [config-loading](config-loading/) | 多环境配置加载 | dev/prod 环境切换 |
| [viper-config](viper-config/) | Viper 配置管理 | 热重载、多格式支持 |
| [config-center](config-center/) | 配置中心集成 | Nacos/Consul/Etcd 配置中心 |

### 安全模块

| 示例 | 说明 | 关键特性 |
|------|------|----------|
| [jwt-example](jwt-example/) | JWT 工具使用示例 | Token 生成与验证 |
| [jwt-security](jwt-security/) | JWT + Security 模块集成 | 认证中间件、权限控制 |

### 缓存模块

| 示例 | 说明 | 关键特性 |
|------|------|----------|
| [cache-example](cache-example/) | 缓存基础操作 | Get/Set/Delete、TTL 管理 |
| [redis-cache](redis-cache/) | Redis 缓存示例 | 连接池、内存回退 |

### 授权模块

| 示例 | 说明 | 关键特性 |
|------|------|----------|
| [casbin-gin](casbin-gin/) | Casbin + Gin 权限控制 | RBAC 模型、路由权限 |
| [casbin-hertz](casbin-hertz/) | Casbin + Hertz 权限控制 | 策略管理、权限校验 |

### 日志模块

| 示例 | 说明 | 关键特性 |
|------|------|----------|
| [zap-logger](zap-logger/) | Zap 日志集成 | 多级别日志、结构化输出 |
| [slog-logger](slog-logger/) | Slog 日志集成 | 标准库 slog 支持 |
| [logrus-logger](logrus-logger/) | Logrus 日志集成 | JSON/Text 格式输出 |

### 可观测性模块

| 示例 | 说明 | 关键特性 |
|------|------|----------|
| [event-bus](event-bus/) | 事件总线示例 | 发布/订阅、异步事件 |
| [health-example](health-example/) | 健康检查端点 | 健康指标、聚合检查 |
| [actuator-example](actuator-example/) | Actuator 监控端点 | 健康、指标、环境信息 |
| [opentelemetry-optimized](opentelemetry-optimized/) | OpenTelemetry 追踪集成 | 分布式追踪、Span 导出 |
| [metrics-example](metrics-example/) | 指标收集示例 | Counter、Gauge、Registry |

### gRPC 模块

| 示例 | 说明 | 关键特性 |
|------|------|----------|
| [grpc](grpc/) | gRPC 完整示例 | 一元 RPC、服务端流式 RPC |
| &nbsp;&nbsp;├─ [grpc-server](grpc/grpc-server/) | gRPC 服务端 | 服务注册、RPC 实现 |
| &nbsp;&nbsp;├─ [grpc-client](grpc/grpc-client/) | gRPC 客户端 | RPC 调用、超时控制 |
| &nbsp;&nbsp;├─ [grpc-server-health](grpc/grpc-server-health/) | gRPC 健康检查 | grpc.health.v1 协议 |
| &nbsp;&nbsp;└─ [interceptor-demo](grpc/interceptor-demo/) | gRPC 拦截器与追踪 | 全局/服务级拦截器 |

### HTTP 客户端

| 示例 | 说明 | 关键特性 |
|------|------|----------|
| [fasthttp-client](fasthttp-client/) | FastHTTP 客户端示例 | 高性能 HTTP 请求 |

### 大型系统示例

| 示例 | 说明 | 架构特点 |
|------|------|----------|
| [e-commerce-system](e-commerce-system/) | 电商系统示例 | 完整业务流程、多模块协同 |
| [blog-system](blog-system/) | 博客系统示例 | CRUD 完整实现、权限管理 |
| [microservice-architecture](microservice-architecture/) | 微服务架构示例 | 服务拆分、网关路由 |
| &nbsp;&nbsp;├─ [api-gateway](microservice-architecture/api-gateway/) | API 网关服务 | 路由转发、负载均衡 |
| &nbsp;&nbsp;├─ [user-service](microservice-architecture/user-service/) | 用户管理服务 | 用户 CRUD、认证 |
| &nbsp;&nbsp;└─ [order-service](microservice-architecture/order-service/) | 订单管理服务 | 订单流程、状态机 |
| [realtime-chat](realtime-chat/) | 实时聊天系统示例 | WebSocket、房间管理、消息广播 |
| [task-management](task-management/) | 任务管理系统示例 | 任务 CRUD、状态跟踪、定时任务 |

---

## 技术特点

所有示例项目都展示了 go-boot 框架的核心特性：

| 特性 | 说明 | 相关示例 |
|------|------|----------|
| **依赖注入容器** | 简化对象创建和依赖管理 | core-* 系列 |
| **AOP 支持** | 实现横切关注点的统一处理 | aop-* 系列 |
| **事件驱动** | 构建松耦合的系统架构 | event-bus |
| **生命周期管理** | 确保资源的正确初始化和清理 | core-integrated |
| **自动配置** | 减少样板代码，提高开发效率 | config-* 系列 |
| **健康检查** | 服务可用性监控 | health-example, actuator-example |
| **指标收集** | 运行时性能监控 | metrics-example |
| **分布式追踪** | 请求链路追踪 | opentelemetry-optimized |
| **gRPC 支持** | 高性能 RPC 通信 | grpc/* 系列 |
| **权限控制** | 基于角色的访问控制 | casbin-* 系列 |
| **JWT 认证** | Token 生成与验证 | jwt-* 系列 |
| **多日志支持** | Zap/Slog/Logrus 集成 | *-logger 系列 |

---

## 学习路径

### 入门级
1. [core-basic](core-basic/) - 了解 Bean 注册与获取
2. [gin-hello](gin-hello/) - 快速启动 HTTP 服务
3. [zap-logger](zap-logger/) - 添加日志支持

### 进阶级
1. [core-injection](core-injection/) - 深入依赖注入
2. [aop-integrated](aop-integrated/) - 掌握切面编程
3. [gin-integrated](gin-integrated/) - Web + 数据库集成
4. [combine-db-cache](combine-db-cache/) - Cache-Aside 缓存模式

### 高级级
1. [microservice-architecture](microservice-architecture/) - 微服务架构设计
2. [e-commerce-system](e-commerce-system/) - 完整业务系统
3. [opentelemetry-optimized](opentelemetry-optimized/) - 可观测性集成
4. [grpc](grpc/) - gRPC 服务与客户端

---

这些示例项目为学习和使用 go-boot 框架提供了实用的参考和指导。