# 字段注入示例

演示 go-boot 核心的字段注入功能，通过 `inject` 结构体标签自动注入依赖。

## 功能特性

- 在结构体字段上使用 `inject` 标签声明依赖
- 容器自动将已注册的 Bean 注入到匹配字段
- 只注入带有 `inject` 标签的字段，普通字段不受影响
- 支持任意类型的依赖注入

## 快速开始

```bash
cd examples/core-injection
go mod tidy
go run .
```

## 预期输出

```
Container created

Before Inject():
  Name: MyService
  DB: <nil> (should be nil)
  Log: <nil> (should be nil)

After Inject():
  Name: MyService (unchanged)
  DB.URL: localhost:5432 (injected!)
  Log.Level: info (injected!)

core-injection example completed successfully!
```

## 代码结构

```go
// 定义结构体
type MyService struct {
    Name string
    DB   *Database `inject:""`
    Log  *Logger   `inject:""`
}

// 注册 Bean
container.Register("database", core.Bean(&Database{URL: "localhost:5432"}))
container.Register("logger", core.Bean(&Logger{Level: "info"}))

// 创建并注入
service := &MyService{Name: "MyService"}
container.Inject(service)
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `inject` 标签 | 声明依赖注入字段 |
| `container.Inject()` | 自动注入依赖 |
| 类型匹配 | 按类型查找并注入 Bean |