// Package main 演示 go-boot 核心 - 工厂函数
//
// 本示例展示如何:
//   - 使用 core.Factory() 工厂函数
//   - 动态创建 Bean
//
// 使用方式:
//
//	cd examples/core-factory
//	go run .
package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/xudefa/go-boot/core"
)

// Config 保存应用配置
type Config struct {
	Env  string
	Port int
}

func main() {
	// 第1步: 创建容器
	container := core.New()
	fmt.Println("Container created")
	fmt.Println()

	// 第2步: 使用工厂函数注册 Bean
	// Factory 接收两个参数：工厂函数和返回类型，容器在首次获取时调用工厂函数
	fmt.Println("Registering factory-created bean...")
	err := container.Register("config",
		core.Factory(
			func(c core.Container) (any, error) {
				fmt.Println("  Factory: creating Config...")
				return &Config{
					Env:  "production",
					Port: 8080,
				}, nil
			},
			reflect.TypeFor[Config](),
		),
		core.Singleton(),
	)
	if err != nil {
		log.Fatalf("Failed to register factory bean: %v", err)
	}

	// 第3步: 获取工厂创建的 Bean
	// 容器会自动调用工厂函数创建实例，之后按 Singleton 作用域缓存
	bean, err := container.Get("config")
	if err != nil {
		log.Fatal(err)
	}
	config := bean.(*Config)
	fmt.Printf("Retrieved Config: Env=%s, Port=%d\n", config.Env, config.Port)

	fmt.Println("\ncore-factory example completed successfully!")
}
