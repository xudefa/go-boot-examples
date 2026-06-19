# 多环境配置加载示例

演示 go-boot 的多环境配置加载功能，支持基础配置与环境特定配置的叠加。

## 功能特性

- **默认配置加载**：自动搜索并加载 `application.json` 配置文件
- **环境特定配置**：支持加载 `application-{profile}.json` 配置文件
- **YAML 配置支持**：支持 YAML 格式的配置文件
- **自定义配置源**：支持用户自定义 PropertySource

## 配置文件搜索路径

按以下优先级搜索配置文件：

1. `/etc/config`
2. 当前目录
3. `./config`
4. 可执行文件所在目录
5. 可执行文件所在目录的 `./config`

## 配置加载顺序

1. 加载基础配置（如 `application.json`）
2. 加载环境特定配置（如 `application-dev.json`）
3. 环境特定配置覆盖基础配置

## 快速开始

```bash
cd examples/config-loading
go mod tidy
go run main.go
```

## 使用示例

### 使用默认配置

```go
app, err := boot.NewApplication(
    boot.WithProfiles("dev"),
)
```

### 使用 YAML 配置

```go
app, err := boot.NewApplication(
    boot.WithConfigType("yaml"),
    boot.WithConfigLocation("./application-prod.yaml"),
)
```

### 使用自定义配置源

```go
customSource := environment.NewMapPropertySource(
    "custom",
    environment.PriorityHigh,
    map[string]any{
        "app.name": "custom-app",
    },
)

app, err := boot.NewApplication(
    boot.WithPropertySource(customSource),
)
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `boot.WithProfiles()` | 指定运行环境 |
| 配置文件搜索 | 多路径自动查找 |
| 配置叠加 | 环境配置覆盖基础配置 |
| 自定义 PropertySource | 灵活的配置来源 |