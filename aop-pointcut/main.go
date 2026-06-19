// Package main 演示 go-boot AOP - 多种切点匹配器
//
// 本示例展示如何:
//   - 使用 MatchByName() 精确匹配方法名
//   - 使用 MatchByNamePrefix() 前缀匹配
//   - 使用 MatchByRegex() 正则匹配
//
// 使用方式:
//
//	cd examples/aop-pointcut
//	go run .
package main

import (
	"fmt"

	"github.com/xudefa/go-boot/aop"
)

// Service 包含多个方法
type Service struct{}

func (s *Service) DoWork() {
	fmt.Println("  [Target] DoWork() called")
}

func (s *Service) DoSomething(name string) {
	fmt.Printf("  [Target] DoSomething(%s) called\n", name)
}

func (s *Service) HandleRequest(id int) {
	fmt.Printf("  [Target] HandleRequest(%d) called\n", id)
}

func (s *Service) ProcessData(data string) {
	fmt.Printf("  [Target] ProcessData(%s) called\n", data)
}

func main() {
	fmt.Println("=== AOP PointCut Matchers Example ===")
	fmt.Println()

	target := &Service{}

	// 示例1: MatchByName
	// MatchByName 精确匹配指定名称的方法，其他方法不受影响
	fmt.Println("1. MatchByName(\"DoWork\"):")
	pointcut1 := aop.MatchByName("DoWork")
	advice1 := aop.Before(func(jp aop.JoinPoint) {
		fmt.Printf("   [Before] %s\n", jp.Signature().Name())
	})
	aspect1 := &aop.AspectMeta{
		PointCut: pointcut1,
		Advice:   advice1,
		Order:    1,
	}

	weaver1 := aop.NewWeaver()
	weaver1.AddAspects(aspect1)
	proxy1 := weaver1.Weave(target)

	if reflectiveProxy, ok := aop.AsReflectiveProxy(proxy1); ok {
		reflectiveProxy.Call("DoWork")              // Should trigger advice
		reflectiveProxy.Call("DoSomething", "test") // Should NOT trigger
	}

	fmt.Println()

	// 示例2: MatchByNamePrefix
	fmt.Println("2. MatchByNamePrefix(\"Do\"):")
	pointcut2 := aop.MatchByNamePrefix("Do")
	advice2 := aop.Before(func(jp aop.JoinPoint) {
		fmt.Printf("   [Before-Do] %s\n", jp.Signature().Name())
	})
	aspect2 := &aop.AspectMeta{
		PointCut: pointcut2,
		Advice:   advice2,
		Order:    1,
	}

	weaver2 := aop.NewWeaver()
	weaver2.AddAspects(aspect2)
	proxy2 := weaver2.Weave(target)

	if reflectiveProxy, ok := aop.AsReflectiveProxy(proxy2); ok {
		reflectiveProxy.Call("DoWork")              // Should trigger
		reflectiveProxy.Call("DoSomething", "test") // Should trigger
		reflectiveProxy.Call("HandleRequest", 1)    // Should NOT trigger
	}

	fmt.Println()

	// 示例3: MatchByRegex
	// MatchByRegex 使用正则表达式匹配方法名，灵活定义拦截范围
	fmt.Println("3. MatchByRegex(\"^Handle.*\"):")
	pointcut3 := aop.MatchByRegex("^Handle.*")
	advice3 := aop.Before(func(jp aop.JoinPoint) {
		fmt.Printf("   [Before-Handle] %s\n", jp.Signature().Name())
	})
	aspect3 := &aop.AspectMeta{
		PointCut: pointcut3,
		Advice:   advice3,
		Order:    1,
	}

	weaver3 := aop.NewWeaver()
	weaver3.AddAspects(aspect3)
	proxy3 := weaver3.Weave(target)

	if reflectiveProxy, ok := aop.AsReflectiveProxy(proxy3); ok {
		reflectiveProxy.Call("DoWork")              // Should NOT trigger
		reflectiveProxy.Call("HandleRequest", 1)    // Should trigger
		reflectiveProxy.Call("ProcessData", "test") // Should NOT trigger
	}

	fmt.Println("\naop-pointcut example completed successfully!")
}
