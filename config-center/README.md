# 配置中心示例

演示如何在 go-boot 应用中使用配置中心功能，支持 Nacos、Etcd、Consul 三种配置中心。

## 功能特性

- 支持 Nacos、Etcd、Consul 三种配置中心
- 在应用启动的配置阶段自动加载远程配置
- 配置中心配置优先级高于本地配置文件
- 支持自定义配置中心参数

## 快速开始

```bash
cd examples/config-center
go mod tidy
go run main.go
```

## 使用示例

### Nacos 配置中心

```go
import (
    "github.com/xudefa/go-boot/boot"
    _ "github.com/xudefa/go-boot/nacos"
)

app, err := boot.NewApplication(
    boot.WithAppName("my-app"),
    boot.WithConfigCenter("nacos", []string{"127.0.0.1:8848"},
        boot.WithConfigCenterDataID("app-config"),
        boot.WithConfigCenterGroup("DEFAULT_GROUP"),
        boot.WithConfigCenterTimeout(5*time.Second),
    ),
)
```

### Etcd 配置中心

```go
import (
    "github.com/xudefa/go-boot/boot"
    _ "github.com/xudefa/go-boot/etcd"
)

app, err := boot.NewApplication(
    boot.WithAppName("my-app"),
    boot.WithConfigCenter("etcd", []string{"127.0.0.1:2379"},
        boot.WithConfigCenterPrefix("/config"),
        boot.WithConfigCenterTimeout(5*time.Second),
    ),
)
```

### Consul 配置中心

```go
import (
    "github.com/xudefa/go-boot/boot"
    _ "github.com/xudefa/go-boot/consul"
)

app, err := boot.NewApplication(
    boot.WithAppName("my-app"),
    boot.WithConfigCenter("consul", []string{"127.0.0.1:8500"},
        boot.WithConfigCenterPrefix("config"),
        boot.WithConfigCenterTimeout(5*time.Second),
    ),
)
```

## 配置参数

### 通用参数

| 参数 | 说明 |
|------|------|
| `WithConfigCenter(centerType, addr, opts...)` | 启用配置中心 |
| `WithConfigCenterTimeout(timeout)` | 设置超时时间 |

### Nacos 特有参数

| 参数 | 默认值 |
|------|--------|
| `WithConfigCenterDataID(dataID)` | `app-config` |
| `WithConfigCenterGroup(group)` | `DEFAULT_GROUP` |

### Etcd/Consul 特有参数

| 参数 | Etcd 默认值 | Consul 默认值 |
|------|-------------|---------------|
| `WithConfigCenterPrefix(prefix)` | `/config` | `config` |

## 配置优先级

配置加载顺序（从高到低）：

1. 命令行参数（`--key=value`）
2. 环境变量（`GO_BOOT_` 前缀）
3. 配置中心配置
4. 本地配置文件

## 注意事项

1. 使用配置中心前需要先导入对应的包（使用 `_` 导入）
2. 配置中心会在应用启动的配置阶段自动加载
3. 如果配置中心连接失败，应用启动会失败
4. 配置中心配置会作为 `PriorityNormal` 优先级的配置源添加到环境中

## 故障排查

### 错误: unsupported config center type

确保已导入对应的配置中心包：

```go
import (
    _ "github.com/xudefa/go-boot/nacos"
    _ "github.com/xudefa/go-boot/etcd"
    _ "github.com/xudefa/go-boot/consul"
)
```

### 错误: config center address is required

确保提供了配置中心地址：

```go
boot.WithConfigCenter("nacos", []string{"127.0.0.1:8848"}, ...)
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| 配置中心集成 | 远程配置管理 |
| 多配置中心支持 | Nacos/Etcd/Consul |
| 配置优先级 | 多层级配置叠加 |
| 自动加载 | 启动时自动拉取远程配置 |