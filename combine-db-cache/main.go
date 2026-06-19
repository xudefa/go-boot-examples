// Package main 演示 go-boot 的 GORM + Redis 双层组合应用（缓存穿透 + 缓存失效模式）
//
// 功能：
//   - GORM 数据库 + Redis 缓存
//   - Cache-Aside 模式（先读缓存 → 未命中读 DB → 回填缓存）
//   - Getter 回调模式
//   - 写操作时缓存失效
//   - 生命周期管理
//
// 使用方式：
//
//	cd examples/combine-db-cache
//	go run .
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/xudefa/go-boot/boot"
	"github.com/xudefa/go-boot/cache"
	"github.com/xudefa/go-boot/environment"
	ggorm "github.com/xudefa/go-boot-gorm"
	"github.com/xudefa/go-boot-redis"
)

// Article 文章模型，保存到数据库
// 使用 gorm tag 定义字段类型和约束
type Article struct {
	ggorm.BaseModel
	Title   string `gorm:"type:varchar(200);not null" json:"title"` // 文章标题
	Content string `gorm:"type:text" json:"content"`                // 文章内容（长文本）
	Author  string `gorm:"type:varchar(100)" json:"author"`         // 作者
	Views   int    `gorm:"default:0" json:"views"`                  // 阅读数，默认 0
}

// ArticleCache 文章缓存层，实现 Cache-Aside 模式
// 组合数据库 Repository 和缓存接口，提供透明的缓存读写
type ArticleCache struct {
	db    *ggorm.Repository[Article]
	cache cache.Cache
}

func NewArticleCache(db *ggorm.Repository[Article], c cache.Cache) *ArticleCache {
	return &ArticleCache{db: db, cache: c}
}

// GetByID 按 ID 获取文章，实现缓存穿透保护
// 优先从缓存读取，未命中时查数据库并回填
func (ac *ArticleCache) GetByID(ctx context.Context, id uint) (*Article, error) {
	// 无缓存时直接查数据库
	if ac.cache == nil {
		return ac.db.FindByID(id)
	}

	// 第一步：查缓存
	cacheKey := fmt.Sprintf("article:%d", id)
	cached, err := ac.cache.Get(ctx, cacheKey)
	if err == nil {
		if article, ok := cached.(*Article); ok {
			return article, nil
		}
	}

	// 第二步：缓存未命中，查数据库
	article, err := ac.db.FindByID(id)
	if err != nil {
		return nil, err
	}
	fmt.Printf("  [DB] 查询文章 ID=%d\n", id)
	// 第三步：回填缓存，设置 10 分钟过期
	_ = ac.cache.Set(ctx, cacheKey, article, 10*time.Minute)
	return article, nil
}

// Create 创建文章，写 DB 后删除缓存
// 采用"删除缓存"而非"更新缓存"策略，避免并发写导致的缓存不一致
func (ac *ArticleCache) Create(ctx context.Context, article *Article) error {
	if err := ac.db.Create(article); err != nil {
		return err
	}
	if ac.cache != nil {
		_ = ac.cache.Del(ctx, fmt.Sprintf("article:%d", article.ID))
	}
	return nil
}

// Update 更新文章，更新 DB 后使缓存失效
func (ac *ArticleCache) Update(ctx context.Context, article *Article) error {
	if err := ac.db.Update(article); err != nil {
		return err
	}
	if ac.cache != nil {
		_ = ac.cache.Del(ctx, fmt.Sprintf("article:%d", article.ID))
	}
	return nil
}

// Delete 删除文章，删除 DB 记录后清除缓存
func (ac *ArticleCache) Delete(ctx context.Context, id uint) error {
	if err := ac.db.Delete(id); err != nil {
		return err
	}
	if ac.cache != nil {
		_ = ac.cache.Del(ctx, fmt.Sprintf("article:%d", id))
	}
	return nil
}

func main() {
	// 创建应用实例，管理生命周期（启动、停止）
	app, err := boot.NewApplication(
		boot.WithAppName("combine-db-cache"),
		boot.WithVersion("1.0.0"),
	)
	if err != nil {
		log.Fatalf("创建应用失败: %v", err)
	}

	// 配置环境属性源，用于集中管理所有配置项
	env := app.Environment()
	env.AddPropertySource(
		environment.NewMapPropertySource("defaults", environment.PriorityLow, map[string]any{
			"boot.datasource.host":     "localhost",
			"boot.datasource.port":     3306,
			"boot.datasource.username": "gate",
			"boot.datasource.password": "123456",
			"boot.datasource.database": "gate",
			"redis.address":            "localhost:6379",
		}),
	)

	// 连接 MySQL 数据库
	// 从环境配置中读取数据库连接参数
	db, err := ggorm.OpenMySQL(
		ggorm.WithHost(env.GetString("boot.datasource.host", "localhost")),
		ggorm.WithPort(env.GetInt("boot.datasource.port", 3306)),
		ggorm.WithUser(env.GetString("boot.datasource.username", "gate")),
		ggorm.WithPassword(env.GetString("boot.datasource.password", "123456")),
		ggorm.WithDBName(env.GetString("boot.datasource.database", "gate")),
		ggorm.WithCharset("utf8mb4"),
		ggorm.WithParseTime(true),
	)
	if err != nil {
		fmt.Printf("数据库不可用: %v\n", err)
		fmt.Println("将以无数据库模式运行演示")
		demoWithoutDB()
		return
	}
	defer db.Close()

	// 自动迁移数据库表结构
	if err := db.DB().AutoMigrate(&Article{}); err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	fmt.Println("数据库迁移完成")

	// 创建泛型 Repository 用于数据库操作
	repo := ggorm.NewRepository[Article](db.DB())

	// 创建 Redis 缓存实例
	var c cache.Cache
	c, err = redis.NewCacheWithConfig(
		[]redis.ClientOption{
			redis.WithAddress(env.GetString("redis.address", "localhost:6379")),
		},
		redis.WithPrefix("article:"),
	)
	if err != nil {
		// Redis 不可用时不阻塞，后续逻辑通过 ac.cache == nil 判断
		fmt.Printf("Redis 不可用: %v，将以无缓存模式运行\n", err)
	}

	// 组合数据库和缓存，创建缓存层
	ac := NewArticleCache(repo, c)

	if err := app.Start(); err != nil {
		log.Fatalf("应用启动失败: %v", err)
	}

	// 执行 Cache-Aside 模式演示
	demoCacheAside(ac)

	defer app.Stop()
}

// demoCacheAside 演示 Cache-Aside 模式的完整流程
// 包括：创建 → 首次查询（缓存未命中）→ 再次查询（命中缓存）→ 更新 → 重新查询
func demoCacheAside(ac *ArticleCache) {
	ctx := context.Background()

	fmt.Println("\n=== Cache-Aside 模式演示 ===")

	article := &Article{
		Title:   "Go-Boot 缓存模式详解",
		Content: "本文详细介绍了 Cache-Aside、Read-Through、Write-Through 等缓存模式...",
		Author:  "张三",
	}
	if err := ac.Create(ctx, article); err != nil {
		fmt.Printf("创建失败: %v\n", err)
		return
	}
	fmt.Printf("创建文章: ID=%d, Title=%s\n", article.ID, article.Title)

	fmt.Println("\n--- 第一次查询（缓存未命中 → 回填）---")
	result, err := ac.GetByID(ctx, article.ID)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	fmt.Printf("结果: Title=%s, Author=%s\n", result.Title, result.Author)

	fmt.Println("\n--- 第二次查询（命中缓存）---")
	result, err = ac.GetByID(ctx, article.ID)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	fmt.Printf("结果: Title=%s, Author=%s\n", result.Title, result.Author)

	fmt.Println("\n--- 更新文章（缓存失效）---")
	article.Title = "Go-Boot 缓存模式（更新版）"
	if err := ac.Update(ctx, article); err != nil {
		fmt.Printf("更新失败: %v\n", err)
		return
	}
	fmt.Println("更新成功，缓存已失效")

	fmt.Println("\n--- 更新后查询（重新加载）---")
	result, err = ac.GetByID(ctx, article.ID)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}
	fmt.Printf("结果: Title=%s, Author=%s\n", result.Title, result.Author)

	fmt.Println("\nGORM + Redis 双层组合示例完成!")
}

// demoWithoutDB 无数据库时的演示模式，仅输出代码结构说明
func demoWithoutDB() {
	fmt.Println("\n演示模式：展示代码结构")
	fmt.Println("实现了 Cache-Aside 模式：")
	fmt.Println("  1. GET: 读缓存 → 未命中 → 读 DB → 回填缓存")
	fmt.Println("  2. CREATE: 写 DB → 删除缓存")
	fmt.Println("  3. UPDATE: 更新 DB → 删除缓存")
	fmt.Println("  4. DELETE: 删除 DB → 删除缓存")
}
