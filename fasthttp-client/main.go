// Package main 演示 go-boot 与 FastHTTP - HTTP 客户端
//
// 本示例展示如何:
//   - 创建 FastHTTP 客户端
//   - 发送 GET/POST 请求
//   - 使用请求选项
//
// 使用方式:
//
//	cd examples/fasthttp-client
//	go run .
package main

import (
	"context"
	"fmt"
	"time"

	fasthttp "github.com/xudefa/go-boot-fasthttp/client"
	"github.com/xudefa/go-boot/net"
)

func main() {
	// 打印示例标题，标识当前运行的是 FastHTTP HTTP 客户端示例
	fmt.Println("=== Fasthttp HTTP Client Example ===")
	fmt.Println()

	// 创建一个后台 Context，用于控制请求的生命周期和取消
	ctx := context.Background()

	// 使用函数式选项模式创建 FastHTTP 客户端
	// WithBaseURL 设置基础 URL，所有相对路径的请求都会拼接该前缀
	// WithTimeout 设置请求超时时间，超过该时间未响应则返回错误
	client, err := fasthttp.NewHttpClient(
		fasthttp.WithBaseURL("http://httpbin.org"),
		fasthttp.WithTimeout(10*time.Second),
	)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}
	defer client.Close()
	fmt.Println("Fasthttp HTTP client created")

	fmt.Println("\nNote: This example requires network access.")
	fmt.Println("The requests below will only work with internet connectivity.")

	// 发送 GET 请求，通过 WithQuery 添加查询参数 name=go-boot
	// 如果网络不可用，请求会失败并打印提示信息
	resp, err := client.Get(ctx, "/get", net.WithQuery("name", "go-boot"))
	if err != nil {
		fmt.Printf("GET request failed (expected without network): %v\n", err)
	} else {
		fmt.Printf("GET response status: %d\n", resp.StatusCode)
		fmt.Printf("Response body: %s\n", string(resp.Body))
	}

	fmt.Println("\nfasthttp-client example completed!")
}
