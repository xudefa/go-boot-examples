# FastHTTP HTTP 客户端

演示 go-boot 基于 FastHTTP 的高性能 HTTP 客户端。

## 功能特性

- 创建 FastHTTP 客户端，设置 Base URL 和超时
- 发送 GET 请求并添加查询参数
- 封装 `net.Client` 统一接口

## 快速开始

```bash
cd examples/fasthttp-client
go mod tidy
go run main.go
```

> 需要网络连接（访问 `httpbin.org`）。无网络时请求会打印错误信息。

## 预期输出

```
=== FastHTTP Client Example ===

Client created with base URL: https://httpbin.org

GET /get?name=go-boot&version=1.0
Status: 200
Response: {...}
```

## 展示的技术点

| 技术点 | 说明 |
|--------|------|
| `fasthttp.NewClient()` | 创建 FastHTTP 客户端 |
| `net.Client` 接口 | 统一的 HTTP 客户端抽象 |
| Base URL | 请求基础地址配置 |
| Timeout | 请求超时控制 |

## 依赖

- [valyala/fasthttp](https://github.com/valyala/fasthttp) — 高性能 HTTP 库
- `github.com/xudefa/go-boot/fasthttp` — FastHTTP 客户端集成
- `github.com/xudefa/go-boot/net` — HTTP 客户端抽象