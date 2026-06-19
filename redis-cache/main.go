// Package main 演示 go-boot 的 Redis 缓存适配
//
// 本示例展示如何:
//   - 创建 Redis 缓存实例
//   - 在 Redis 不可用时回退到内存缓存
//   - 使用 Cache 接口的统一操作方法
//
// 使用方式:
//
//	cd examples/redis-cache
//	go run .
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xudefa/go-boot/cache"

	rboot "github.com/xudefa/go-boot-redis"
)

func main() {
	fmt.Println("=== Redis Cache Example ===")
	fmt.Println()

	// 创建背景上下文，用于传递请求生命周期信息
	ctx := context.Background()

	// 配置 Redis 客户端连接参数
	clientOpts := []rboot.ClientOption{
		rboot.WithAddress("localhost:6379"),
		rboot.WithDB(0),
	}

	// 创建缓存实例时同时指定客户端配置和缓存专用选项
	// NewCacheWithConfig 封装了 Redis 客户端创建和 Cache 接口包装
	myCache, err := rboot.NewCacheWithConfig(clientOpts,
		rboot.WithPrefix("myapp:"),           // 缓存键前缀，避免与其它应用冲突
		rboot.WithDefaultTTL(10*time.Second), // 默认过期时间
	)
	if err != nil {
		// Redis 不可用时自动降级为内存缓存，保证演示能继续
		fmt.Printf("NewCache failed: %v\n", err)
		fmt.Println("Creating local-only demo instead...")

		localCache := cache.NewMemoryCache()
		demoRedisCache(ctx, localCache)
		return
	}
	defer myCache.Close()

	fmt.Println("Redis cache created successfully!")
	demoRedisCacheWithRedis(ctx, myCache)
}

// demoRedisCache 使用统一的 cache.Cache 接口操作缓存
// 即使底层是内存缓存，操作接口完全一致
func demoRedisCache(ctx context.Context, c cache.Cache) {
	fmt.Println("Demo with MemoryCache (Redis unavailable):")
	c.Set(ctx, "user:1", "Alice", 5*time.Minute)
	c.Set(ctx, "user:2", "Bob", 5*time.Minute)

	val, _ := c.Get(ctx, "user:1")
	fmt.Printf("  Get user:1 = %v\n", val)

	exists, _ := c.Exists(ctx, "user:1")
	fmt.Printf("  Exists user:1 = %v\n", exists)

	c.Del(ctx, "user:2")
	fmt.Println("  Deleted user:2")

	// 通过类型断言检查是否支持 Clear 方法（MemoryCache 支持）
	if cv, ok := c.(interface{ Clear(context.Context) error }); ok {
		cv.Clear(ctx)
		fmt.Println("  Cache cleared")
	}
	fmt.Println("\nredis-cache example completed!")
}

// demoRedisCacheWithRedis 演示 Redis 可用时的操作
func demoRedisCacheWithRedis(ctx context.Context, c cache.Cache) {
	defer c.Close()

	fmt.Println("Redis cache demo:")

	// 存储字符串值，设置 1 分钟过期
	c.Set(ctx, "greeting", "Hello from Redis!", time.Minute)
	fmt.Println("  Set greeting -> 'Hello from Redis!'")

	// 读取缓存值
	val, err := c.Get(ctx, "greeting")
	if err != nil {
		fmt.Printf("  Get error: %v\n", err)
	} else {
		fmt.Printf("  Get greeting = %v\n", val)
	}

	// 存储复杂类型（map）
	user := map[string]any{
		"id":    1,
		"name":  "Alice",
		"email": "alice@example.com",
	}
	c.Set(ctx, "user:1", user, 5*time.Minute)
	fmt.Println("  Set user:1 -> map")

	// 获取复杂类型
	cachedUser, err := c.Get(ctx, "user:1")
	if err != nil {
		fmt.Printf("  Get user:1 error: %v\n", err)
	} else {
		fmt.Printf("  Get user:1 = %v\n", cachedUser)
	}

	// 检查键是否存在
	exists, _ := c.Exists(ctx, "greeting")
	fmt.Printf("  Exists greeting: %v\n", exists)

	// 删除键
	c.Del(ctx, "greeting")
	fmt.Println("  Del greeting")

	exists, _ = c.Exists(ctx, "greeting")
	fmt.Printf("  Exists greeting after Del: %v\n", exists)

	fmt.Println("\nRedis cache demo completed!")
}
