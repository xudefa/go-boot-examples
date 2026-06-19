# 属性源示例

演示 go-boot 框架的属性源（PropertySource）功能，支持多来源配置和优先级管理。

## 功能特性

- **MapPropertySource**：基于 map 的属性源实现
- **优先级管理**：支持不同优先级的配置源叠加
- **PropertySource 接口**：统一的属性源抽象

## 快速开始

```bash
cd examples/property-source
go mod tidy
go run main.go
```

## 使用示例

### 创建 MapPropertySource

```go
source := environment.NewMapPropertySource(
    "custom",
    environment.PriorityHigh,
    map[string]any{
        "app.name": "my-app",
        "server.port": 8080,
    },
)
```

### 添加到环境

```go
env.AddPropertySource(source)
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| MapPropertySource | 基于 map 的属性源 |
| 优先级管理 | 配置源叠加机制 |
| PropertySource 接口 | 统一的属性源抽象 |