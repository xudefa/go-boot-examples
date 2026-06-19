# Viper 配置管理示例

演示 go-boot 与 Viper 集成的配置管理功能，支持多源配置、默认值和上下文感知。

## 功能特性

- **Viper 配置管理器**：通过 `viper.New()` 创建，支持多源配置
- **FileLoader**：文件加载器，支持配置文件名、路径、环境感知
- **MustNew 默认值**：提供默认值机制（`SetDefault` / `SetDefaults`）
- **上下文创建**：通过 `NewWithContext` 支持上下文感知的配置创建

## 快速开始

```bash
cd examples/viper-config
go mod tidy
go run main.go
```

## 使用示例

### 创建 Viper 配置

```go
cfg := viper.New()
```

### 设置默认值

```go
cfg.SetDefault("app.name", "my-app")
cfg.SetDefault("server.port", 8080)
```

### 文件加载器

```go
loader := viper.NewFileLoader(
    viper.WithConfigName("application"),
    viper.WithConfigPath("./config"),
    viper.WithConfigType("json"),
)
```

### 上下文感知创建

```go
cfg, err := viper.NewWithContext(ctx,
    viper.WithDefault("app.name", "my-app"),
)
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `viper.New()` | 创建 Viper 配置管理器 |
| FileLoader | 文件配置加载 |
| SetDefault | 默认值机制 |
| NewWithContext | 上下文感知创建 |