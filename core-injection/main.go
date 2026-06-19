// Package main 演示 go-boot 核心 - 字段注入与 inject 标签
//
// 本示例展示如何:
//   - 定义带 inject 标签的结构体
//   - 使用 container.Inject() 填充字段
//
// 使用方式:
//
//	cd examples/core-injection
//	go run .
package main

import (
	"fmt"
	"log"

	"github.com/xudefa/go-boot/core"
)

// Database 模拟数据库连接
type Database struct {
	URL string
}

// Logger 模拟日志记录器
type Logger struct {
	Level string
}

// UserService 演示字段注入
type UserService struct {
	// inject 标签告诉容器将 "db" Bean 注入到此字段
	DB *Database `inject:"db"`

	// logger 的 inject 标签
	Log *Logger `inject:"logger"`

	// 没有 inject 标签的字段不会被注入
	Name string
}

func main() {
	// 第1步: 创建容器
	container := core.New()
	fmt.Println("Container created")

	// 第2步: 注册依赖
	// 先注册依赖 Bean，后续注入时容器会按 ID 查找
	err := container.Register("db", core.Bean(&Database{URL: "localhost:5432"}))
	if err != nil {
		log.Fatalf("Failed to register db: %v", err)
	}

	err = container.Register("logger", core.Bean(&Logger{Level: "info"}))
	if err != nil {
		log.Fatalf("Failed to register logger: %v", err)
	}

	// 第3步: 注册服务（无需手动注入）
	// 注意：当使用 core.Bean() 注册时，Get() 返回的是同一个实例，
	// inject 标签会在 createInstance 时自动注入。
	// 如果要演示"注入前"的状态，需要先获取原始实例再手动注入。
	// 这里我们直接注册一个不带 inject 标签的 UserService 来演示差异。
	err = container.Register("userService", core.Bean(&UserService{Name: "MyService"}))
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// 第4步: 获取服务（依赖已自动注入）
	// Get() 内部的 createInstance 会自动处理 inject 标签
	bean, err := container.Get("userService")
	if err != nil {
		log.Fatalf("Failed to get bean: %v", err)
	}

	service := bean.(*UserService)
	fmt.Printf("\nAfter Get() (auto-injected):\n")
	fmt.Printf("  Name: %s\n", service.Name)
	if service.DB != nil {
		fmt.Printf("  DB.URL: %s (injected!)\n", service.DB.URL)
	} else {
		fmt.Printf("  DB: nil (not injected)\n")
	}
	if service.Log != nil {
		fmt.Printf("  Log.Level: %s (injected!)\n", service.Log.Level)
	} else {
		fmt.Printf("  Log: nil (not injected)\n")
	}

	// 第5步: 演示手动 Inject
	// 创建一个新实例并手动注入，展示 Inject() 的用法
	manualService := &UserService{Name: "ManualService"}
	fmt.Printf("\nBefore Inject():\n")
	fmt.Printf("  Name: %s\n", manualService.Name)
	fmt.Printf("  DB: %v (should be nil)\n", manualService.DB)
	fmt.Printf("  Log: %v (should be nil)\n", manualService.Log)

	err = container.Inject(manualService)
	if err != nil {
		log.Fatalf("Failed to inject: %v", err)
	}

	fmt.Printf("\nAfter Inject():\n")
	fmt.Printf("  Name: %s (unchanged)\n", manualService.Name)
	if manualService.DB != nil {
		fmt.Printf("  DB.URL: %s (injected!)\n", manualService.DB.URL)
	}
	if manualService.Log != nil {
		fmt.Printf("  Log.Level: %s (injected!)\n", manualService.Log.Level)
	}

	fmt.Println("\ncore-injection example completed successfully!")
}
