// Package main 演示 go-boot 核心 - 基础 Bean 注册与获取
//
// 本示例展示如何:
//   - 创建容器
//   - 使用 core.Bean() 注册 Bean
//   - 使用 container.Get() 获取 Bean
//
// 使用方式:
//
//	cd examples/core-basic
//	go run .
package main

import (
	"fmt"
	"log"

	"github.com/xudefa/go-boot/core"
)

// UserService 是一个简单的服务结构体
type UserService struct {
	Name string
}

func main() {
	// 第1步: 创建新容器
	// core.New() 返回一个默认配置的 IoC 容器实例
	container := core.New()
	fmt.Println("Container created")

	// 第2步: 使用 core.Bean() 注册 Bean
	// 使用 ID "userService" 注册 UserService 实例
	err := container.Register("userService", core.Bean(&UserService{Name: "MyService"}))
	if err != nil {
		log.Fatalf("Failed to register bean: %v", err)
	}
	fmt.Println("Bean registered: userService")

	// 第3步: 使用 container.Get() 获取 Bean
	// 通过注册时指定的 ID 从容器中获取 Bean
	bean, err := container.Get("userService")
	if err != nil {
		log.Fatalf("Failed to get bean: %v", err)
	}

	// 第4步: 类型断言并使用 Bean
	service, ok := bean.(*UserService)
	if !ok {
		log.Fatal("Bean is not of type *UserService")
	}
	fmt.Printf("Retrieved bean: %+v\n", service)

	// 第5步: 演示 Has() 方法
	// Has() 判断容器中是否存在指定 ID 的 Bean，返回布尔值
	if container.Has("userService") {
		fmt.Println("Container has userService: true")
	}

	fmt.Println("\ncore-basic example completed successfully!")
}
