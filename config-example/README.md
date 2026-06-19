# 配置管理示例

演示 go-boot 框架的配置管理抽象层，包括配置加载、校验和变更监听。

## 功能特性

- **Config 接口**：展示 `Get`、`GetString`、`GetInt`、`Unmarshal` 等核心方法
- **ConfigOption**：函数式选项模式配置（名称、路径、类型、环境）
- **LoaderChain**：加载器链，支持多来源配置加载
- **Validator**：配置校验器，支持必填校验、范围校验、正则校验、枚举校验
- **WatchManager**：配置变更监听，支持 Modify/Delete/Create 事件

## 快速开始

```bash
cd examples/config-example
go mod tidy
go run main.go
```

## 使用示例

### 基础配置获取

```go
val := config.GetString("app.name")
port := config.GetInt("server.port")
```

### 配置校验

```go
validator := config.NewValidator()
validator.Required("app.name")
validator.Range("server.port", 1024, 65535)
```

### 配置变更监听

```go
watchManager.OnChange(func(event config.WatchEvent) {
    fmt.Printf("Config changed: %s\n", event.Key)
})
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| Config 接口 | 统一的配置访问抽象 |
| 函数式选项 | 灵活的配置创建方式 |
| 加载器链 | 多来源配置加载 |
| 配置校验 | 确保配置有效性 |
| 变更监听 | 配置热更新支持 |