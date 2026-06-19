// Package main 演示 go-boot 与 Gin + GORM 的集成
//
// 本示例展示如何:
//   - 创建 Gin HTTP 服务器
//   - 连接 MySQL 数据库
//   - 实现完整的 RESTful CRUD API
//
// 使用方式:
//
//	cd examples/gin-integrated
//	go run .
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/xudefa/go-boot/boot"
	"github.com/xudefa/go-boot/core"
	ggin "github.com/xudefa/go-boot-gin/server"
	ggorm "github.com/xudefa/go-boot-gorm"
)

// User 用户模型，对应数据库中的 users 表
// 内嵌 ggorm.BaseModel 自动获得 ID、CreatedAt、UpdatedAt、DeletedAt 字段
// 使用 gorm tag 定义数据库列的约束，json tag 定义序列化字段名
type User struct {
	ggorm.BaseModel
	Name  string `gorm:"type:varchar(100);not null" json:"name"`
	Email string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Age   int    `gorm:"default:0" json:"age"`
}

func main() {
	// 打印示例标题，标识当前运行的是 Gin + GORM 集成示例
	fmt.Println("=== Gin + GORM Integrated Example ===")
	fmt.Println()

	// 使用 boot.NewApplication 创建应用实例
	// 内部包含 Container（IoC 容器）、Environment（环境配置）、Lifecycle（生命周期）等核心组件
	app, err := boot.NewApplication(
		boot.WithAppName("gin-gorm-demo"),
		boot.WithVersion("1.0.0"),
	)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}
	fmt.Println("Application created")

	// 使用函数式选项模式打开 MySQL 数据库连接
	// 如果数据库不可用（比如本地未安装 MySQL），会进入演示模式
	db, err := ggorm.OpenMySQL(
		ggorm.WithHost("localhost"),
		ggorm.WithPort(3306),
		ggorm.WithUser("gate"),
		ggorm.WithPassword("123456"),
		ggorm.WithDBName("gate"),
		ggorm.WithCharset("utf8mb4"),
		ggorm.WithParseTime(true),
	)
	if err != nil {
		fmt.Printf("Cannot connect to MySQL (expected if no DB): %v\n", err)
		fmt.Println("Demo mode: showing API structure without DB connection")
		demoIntegratedAPI()
		return
	}
	defer db.Close()
	fmt.Println("Connected to MySQL")

	// 自动迁移数据库表结构，根据 User 模型自动创建或更新表
	db.DB().AutoMigrate(&User{})
	// 创建泛型 Repository，提供类型安全的 CRUD 操作
	repo := ggorm.NewRepository[User](db.DB())

	// 创建 Gin 服务器，从应用中获取容器
	server := ggin.New(
		ggin.WithContainer(app.Container()),
		ggin.WithMode("debug"),
		ggin.WithHost("localhost"),
		ggin.WithPort(8082),
	)

	// GET /api/users - 查询所有用户列表
	server.GET("/api/users", func(c *gin.Context) {
		users, err := repo.FindAll(nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	// GET /api/users/:id - 根据 ID 查询单个用户
	// 从 URL 路径参数中解析 ID，转换为 uint 类型
	server.GET("/api/users/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		user, err := repo.FindByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	// POST /api/users - 创建新用户
	// 从请求体中绑定 JSON 到 User 结构体，然后写入数据库
	server.POST("/api/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := repo.Create(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	})

	// DELETE /api/users/:id - 根据 ID 删除用户
	// 使用逻辑删除（软删除），实际设置 deleted_at 字段
	server.DELETE("/api/users/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if err := repo.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	})

	fmt.Println("Routes registered:")
	fmt.Println("  GET    /api/users")
	fmt.Println("  GET    /api/users/:id")
	fmt.Println("  POST   /api/users")
	fmt.Println("  DELETE /api/users/:id")

	fmt.Println("\nStarting server on :8082...")
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func demoIntegratedAPI() {
	_ = core.New()
	fmt.Println("\nIntegrated API structure ready.")
	fmt.Println("Routes would be:")
	fmt.Println("  GET    /api/users")
	fmt.Println("  GET    /api/users/:id")
	fmt.Println("  POST   /api/users")
	fmt.Println("  DELETE /api/users/:id")
	fmt.Println("\ngin-integrated example completed!")
}
