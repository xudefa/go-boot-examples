// Package main 演示 go-boot 核心 - 单例与原型作用域
//
// 本示例展示如何:
//   - 注册单例 Bean（默认）
//   - 注册原型 Bean
//   - 理解两者行为的区别
//
// 使用方式:
//
//	cd examples/core-scope
//	go run .
package main

import (
	"fmt"
	"log"

	"github.com/xudefa/go-boot/core"
)

// Counter 是一个简单的计数器结构体
type Counter struct {
	Value int
}

func main() {
	// 第1步: 创建容器
	container := core.New()
	fmt.Println("Container created")
	fmt.Println()

	// 第2步: 注册单例 Bean（默认作用域）
	err := container.Register("singletonCounter",
		core.Bean(&Counter{Value: 0}),
		core.Singleton(), // 这是默认的，但显式声明更清晰
	)
	if err != nil {
		log.Fatalf("Failed to register singleton: %v", err)
	}
	fmt.Println("Registered singleton bean: singletonCounter")

	// 第3步: 多次获取单例 Bean
	// 单例作用域下，多次 Get 返回同一个实例，修改一处将影响所有引用
	c1, err := container.Get("singletonCounter")
	if err != nil {
		log.Fatal(err)
	}
	c1.(*Counter).Value = 100

	c2, err := container.Get("singletonCounter")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Singleton bean behavior:")
	fmt.Printf("  First get:  value = %d\n", c1.(*Counter).Value)
	fmt.Printf("  Second get: value = %d (same instance, same value)\n", c2.(*Counter).Value)
	fmt.Printf("  Same instance: %v\n", c1 == c2)
	fmt.Println()

	// 第4步: 注册原型 Bean
	err = container.Register("prototypeCounter",
		core.Bean(&Counter{Value: 0}),
		core.Prototype(), // New instance each time
	)
	if err != nil {
		log.Fatalf("Failed to register prototype: %v", err)
	}
	fmt.Println("Registered prototype bean: prototypeCounter")

	// 第5步: 多次获取原型 Bean
	// 原型作用域下，每次 Get 都创建新实例，互不影响
	p1, err := container.Get("prototypeCounter")
	if err != nil {
		log.Fatal(err)
	}

	p2, err := container.Get("prototypeCounter")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Prototype bean behavior:")
	fmt.Printf("  First get:  value = %d\n", p1.(*Counter).Value)
	fmt.Printf("  Second get: value = %d\n", p2.(*Counter).Value)
	fmt.Printf("  Same instance: %v (different instances)\n", p1 == p2)
	fmt.Println()

	// 第6步: 独立修改原型实例
	p1.(*Counter).Value = 500
	fmt.Printf("After modifying first instance:\n")
	fmt.Printf("  First instance:  %d\n", p1.(*Counter).Value)
	fmt.Printf("  Second instance: %d (unchanged)\n", p2.(*Counter).Value)

	fmt.Println("\ncore-scope example completed successfully!")
}
