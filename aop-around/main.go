// Package main 演示 go-boot AOP - Around 通知与 proceed
//
// 本示例展示如何:
//   - 创建 Around 通知
//   - 使用 ProceedFunc 控制方法执行
//   - 修改返回值
//
// 使用方式:
//
//	cd examples/aop-around
//	go run .
package main

import (
	"fmt"
	"time"

	"github.com/xudefa/go-boot/aop"
)

// Calculator 目标结构体
type Calculator struct{}

// Add 模拟一个将被拦截的方法
func (c *Calculator) Add(a, b int) int {
	fmt.Printf("  [Target] Add(%d, %d) called\n", a, b)
	time.Sleep(50 * time.Millisecond) // Simulate work
	return a + b
}

func main() {
	fmt.Println("=== AOP Around Advice Example ===")
	fmt.Println()

	// 第1步: 创建目标对象
	target := &Calculator{}

	// 第2步: 创建切点
	pointcut := aop.MatchByName("Add")
	fmt.Println("Created pointcut: MatchByName(\"Add\")")

	// 第3步: 创建带计时的 Around 通知
	// Around 通知接收两个参数：连接点信息和 proceed 回调，调用 proceed 才会执行目标方法
	aroundAdvice := aop.Around(func(jp aop.JoinPoint, proceed aop.ProceedFunc) any {
		start := time.Now()

		fmt.Printf("  [Around] Before: %s, args=%v\n", jp.Signature().Name(), jp.Args())

		// 调用目标方法
		result := proceed(jp.Args()...)

		elapsed := time.Since(start)
		fmt.Printf("  [Around] After: result=%v, elapsed=%v\n", result, elapsed)

		// 可以修改结果
		if intResult, ok := result.(int); ok {
			return intResult * 2 // 将结果翻倍！
		}
		return result
	})
	fmt.Println("Created Around advice (with timing and result modification)")

	// 第4步: 创建切面并织入
	aspect := &aop.AspectMeta{
		PointCut: pointcut,
		Advice:   aroundAdvice,
		Order:    1,
	}

	weaver := aop.NewWeaver()
	weaver.AddAspects(aspect)

	proxy := weaver.Weave(target)
	fmt.Println("Weaved target object")
	fmt.Println()

	// 第5步: 在代理上调用方法
	// 调用代理对象的 Call 方法时会触发 Around 通知，执行前后增强逻辑
	fmt.Println("Calling Add(10, 20) on proxy:")
	if reflectiveProxy, ok := aop.AsReflectiveProxy(proxy); ok {
		result, err := reflectiveProxy.Call("Add", 10, 20)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Final result: %d (note: result was doubled by Around advice)\n", result)
		}
	}

	fmt.Println("\naop-around example completed successfully!")
}
