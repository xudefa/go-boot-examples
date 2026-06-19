// Package main 演示 go-boot AOP - 多通知与排序
//
// 本示例展示如何:
//   - 使用多个通知（Before、After、Around）
//   - 通过 Order 字段控制执行顺序
//   - 理解通知链执行
//
// 使用方式:
//
//	cd examples/aop-integrated
//	go run .
package main

import (
	"fmt"
	"time"

	"github.com/xudefa/go-boot/aop"
)

// OrderService 目标结构体
type OrderService struct{}

func (s *OrderService) CreateOrder(customerID int, amount float64) string {
	fmt.Println("  [Target] CreateOrder called")
	time.Sleep(30 * time.Millisecond)
	return fmt.Sprintf("Order-%d-%.0f", customerID, amount)
}

func main() {
	fmt.Println("=== AOP Integrated Example (Multiple Advices) ===")
	fmt.Println()

	target := &OrderService{}

	// 创建多个不同顺序的通知
	// Order 值越小优先级越高，决定通知在拦截链中的执行位置

	// Before 通知（顺序: 1 - 最先执行）
	beforeAdvice := aop.Before(func(jp aop.JoinPoint) {
		fmt.Printf("[Before-1] Method: %s, Args: %v\n", jp.Signature().Name(), jp.Args())
	})

	// Around 通知（顺序: 2 - 包装方法）
	aroundAdvice := aop.Around(func(jp aop.JoinPoint, proceed aop.ProceedFunc) any {
		start := time.Now()
		fmt.Printf("[Around-2] Before proceed: %s\n", jp.Signature().Name())

		result := proceed(jp.Args()...)

		elapsed := time.Since(start)
		fmt.Printf("[Around-2] After proceed: result=%v, elapsed=%v\n", result, elapsed)
		return result
	})

	// After 通知（顺序: 3 - 方法执行后）
	afterAdvice := aop.After(func(jp aop.JoinPoint) {
		fmt.Printf("[After-3] Method completed: %s\n", jp.Signature().Name())
	})

	// AfterReturning 通知（顺序: 4）
	afterReturningAdvice := aop.AfterReturning(func(jp aop.JoinPoint, result any) {
		fmt.Printf("[AfterReturning-4] Method returned: %v\n", result)
	})

	// 创建切面
	pointcut := aop.MatchByName("CreateOrder")

	aspects := []*aop.AspectMeta{
		{
			PointCut: pointcut,
			Advice:   beforeAdvice,
			Order:    1,
		},
		{
			PointCut: pointcut,
			Advice:   aroundAdvice,
			Order:    2,
		},
		{
			PointCut: pointcut,
			Advice:   afterAdvice,
			Order:    3,
		},
		{
			PointCut: pointcut,
			Advice:   afterReturningAdvice,
			Order:    4,
		},
	}

	// 创建织入器并添加所有切面
	// 添加多个切面时，weaver 会自动按 Order 排序后织入
	weaver := aop.NewWeaver()
	weaver.AddAspects(aspects...)

	proxy := weaver.Weave(target)
	fmt.Println("Created proxy with 4 advices (orders 1-4)")
	fmt.Println()

	// 调用方法
	fmt.Println("Calling CreateOrder(100, 50.0):")
	if reflectiveProxy, ok := aop.AsReflectiveProxy(proxy); ok {
		result, err := reflectiveProxy.Call("CreateOrder", 100, 50.0)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("\nFinal result: %s\n", result)
		}
	}

	fmt.Println("\naop-integrated example completed successfully!")
}
