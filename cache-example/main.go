// Package main 演示 go-boot 的缓存抽象
//
// 本示例展示如何:
//   - 创建内存缓存
//   - 使用 Get/Set/Exists/TTL 等基础操作
//   - 使用批量操作（GetMulti/SetMulti）
//   - 使用 GetWithGetter 延迟加载
//
// 使用方式:
//
//	cd examples/cache-example
//	go run .
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xudefa/go-boot/cache"
)

func main() {
	fmt.Println("=== Cache Example ===")
	fmt.Println()

	// 创建背景上下文，用于传递超时和取消信号
	ctx := context.Background()
	// 创建内存缓存实例，存储在内存中，无需外部依赖
	mc := cache.NewMemoryCache()
	fmt.Println("Memory cache created")

	// 设置缓存项，支持任意类型的 value
	// 第三个参数为 TTL（过期时间）
	mc.Set(ctx, "user:1", "Alice", 5*time.Minute)
	mc.Set(ctx, "user:2", "Bob", 5*time.Minute)
	// 支持存储复杂类型，如 map
	mc.Set(ctx, "session:abc", map[string]any{"token": "xyz"}, 1*time.Hour)
	fmt.Println("Items set: user:1, user:2, session:abc")

	// 使用 Get 方法获取缓存值
	val, err := mc.Get(ctx, "user:1")
	if err != nil {
		fmt.Printf("Get error: %v\n", err)
	} else {
		fmt.Printf("Get user:1 = %v\n", val)
	}

	// 使用 Exists 方法检查键是否存在
	exists, _ := mc.Exists(ctx, "user:2")
	missing, _ := mc.Exists(ctx, "nonexistent")
	fmt.Printf("user:2 exists: %v, nonexistent exists: %v\n", exists, missing)

	// 使用 TTL 方法获取剩余过期时间
	ttl, _ := mc.TTL(ctx, "user:1")
	fmt.Printf("TTL for user:1: %v\n", ttl)

	fmt.Println("\n=== GetMulti / SetMulti ===")
	// 批量获取，未命中的键不会出现在结果中
	items, _ := mc.GetMulti(ctx, []string{"user:1", "user:2", "nonexistent"})
	fmt.Printf("Multi-get results: %+v\n", items)

	// 批量设置，所有键使用相同的 TTL
	mc.SetMulti(ctx, map[string]any{"user:3": "Charlie", "user:4": "Diana"}, time.Minute)
	val3, _ := mc.Get(ctx, "user:3")
	fmt.Printf("After SetMulti, user:3 = %v\n", val3)

	fmt.Println("\n=== GetWithGetter ===")
	// GetWithGetter 实现缓存穿透保护：
	// 缓存未命中时自动调用 Getter 回调加载数据并回填
	getter := cache.Getter(func(ctx context.Context, key string) (any, error) {
		fmt.Printf("  Getter called for key: %s\n", key)
		return "computed-value", nil
	})
	gv, err := mc.GetWithGetter(ctx, "computed:1", getter)
	if err != nil {
		fmt.Printf("GetWithGetter error: %v\n", err)
	} else {
		fmt.Printf("GetWithGetter result: %v\n", gv)
	}

	fmt.Println("\n=== Delete operations ===")
	// 批量删除
	_ = mc.DeleteMulti(ctx, []string{"user:2"})
	e1, _ := mc.Exists(ctx, "user:2")
	fmt.Printf("After DeleteMulti, user:2 exists: %v\n", e1)

	// 清空所有缓存
	_ = mc.Clear(ctx)
	e2, _ := mc.Exists(ctx, "user:1")
	fmt.Printf("After Clear, user:1 exists: %v\n", e2)

	fmt.Println("\ncache example completed successfully!")
}
