// Package main 演示 go-boot AOP - 基础 Before 通知
//
// 本示例展示如何:
//   - 创建简单的 Before 通知
//   - 将通知应用到目标方法
//   - 理解方法拦截
//
// 使用方式:
//
//	cd examples/aop-basic
//	go run .
package main

import (
	"fmt"

	"github.com/xudefa/go-boot/aop"
)

// UserService 目标结构体
type UserService struct {
	Name string
}

// GetUser 模拟一个将被拦截的方法
func (s *UserService) GetUser(id int) string {
	fmt.Printf("  [Target] GetUser(%d) called\n", id)
	return fmt.Sprintf("User%d", id)
}

func main() {
	fmt.Println("=== AOP Basic Example (Before Advice) ===")
	fmt.Println()

	// 第1步: 创建目标对象
	// 目标对象是需要被 AOP 增强的原始对象
	target := &UserService{Name: "MyService"}

	// 第2步: 创建切点以匹配 GetUser 方法
	pointcut := aop.MatchByName("GetUser")
	fmt.Println("Created pointcut: MatchByName(\"GetUser\")")

	// 第3步: 创建 Before 通知
	beforeAdvice := aop.Before(func(jp aop.JoinPoint) {
		fmt.Printf("  [Before] Method: %s\n", jp.Signature().Name())
		fmt.Printf("  [Before] Args: %v\n", jp.Args())
	})
	fmt.Println("Created Before advice")

	// 第4步: 创建 AspectMeta
	// AspectMeta 将切点和通知组合在一起，Order 用于多切面排序
	aspect := &aop.AspectMeta{
		PointCut: pointcut,
		Advice:   beforeAdvice,
		Order:    1,
	}
	fmt.Println("Created AspectMeta")
	fmt.Println()

	// 第5步: 创建织入器并添加切面
	weaver := aop.NewWeaver()
	weaver.AddAspects(aspect)
	fmt.Println("Created weaver and added aspect")
	fmt.Println()

	// 第6步: 织入目标对象
	proxy := weaver.Weave(target)
	fmt.Println("Weaved target object, got proxy")
	fmt.Println()

	// 第7步: 在代理上调用方法（将执行通知）
	// 由于 Go 运行时无法动态替换结构体方法，Weave 返回 ReflectiveAopProxy，
	// 需要通过 Call 方法调用，AOP 通知会自动执行
	fmt.Println("Calling method on proxy:")
	if reflectiveProxy, ok := aop.AsReflectiveProxy(proxy); ok {
		result, err := reflectiveProxy.Call("GetUser", 123)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Result: %s\n", result)
		}
	} else {
		fmt.Println("Failed to get reflective proxy")
	}

	fmt.Println("\naop-basic example completed successfully!")
}
