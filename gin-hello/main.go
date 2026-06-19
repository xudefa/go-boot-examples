// Package main 演示 go-boot 与 Gin - 基础 HTTP 服务器
//
// 本示例展示如何:
//   - 创建 Gin 服务器
//   - 注册路由
//   - 与核心容器集成
//
// 使用方式:
//
//	cd examples/gin-hello
//	go run .
//	# 访问 http://localhost:8080/hello
package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	ggin "github.com/xudefa/go-boot-gin/server"
	"github.com/xudefa/go-boot/core"
)

func main() {
	// 打印示例标题，标识当前运行的是 Gin Hello World 示例
	fmt.Println("=== Gin Hello World Example ===")
	fmt.Println()

	// 创建一个新的 IoC 容器实例，用于管理 Bean 的注册与依赖注入
	container := core.New()
	fmt.Println("Container created")

	// 使用函数式选项模式创建 Gin 服务器
	// WithContainer 将容器注入服务器以便后续扩展
	// WithMode 设置 Gin 的运行模式（debug/release/test）
	// WithHost 设置监听地址和端口
	server := ggin.New(
		ggin.WithContainer(container),
		ggin.WithMode("debug"),
		ggin.WithHost("localhost"),
		ggin.WithPort(8080),
	)
	fmt.Println("Gin server created")

	// 注册一个 GET /hello 路由，返回简单的欢迎信息
	server.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello from go-boot with Gin!")
	})
	fmt.Println("Routes registered: GET /hello")

	// 将 Gin 服务器注册到 IoC 容器中
	// 这样其他组件可以通过容器获取服务器实例
	err := container.Register("ginServer", core.Bean(server))
	if err != nil {
		log.Printf("Warning: Failed to register server: %v\n", err)
	}

	fmt.Println("\nStarting server on :8080...")
	fmt.Println("Visit: http://localhost:8080/hello")
	fmt.Println("Press Ctrl+C to stop")

	// 启动 HTTP 服务器，开始监听并处理请求
	// Start 是阻塞调用，会一直运行直到服务器出错或被中断
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// 测试示例:
//   curl http://localhost:8080/hello
