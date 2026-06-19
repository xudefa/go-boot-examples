// Package main 演示 go-boot 与 GORM 的基础使用
//
// 本示例展示如何:
//   - 连接 MySQL 数据库
//   - 了解可用数据库类型
//   - 了解 Option 函数和 Open 函数
//   - 了解 Repository[T] 方法
//
// 使用方式:
//
//	cd examples/gorm-basic
//	go run .
package main

import (
	"fmt"

	"github.com/xudefa/go-boot-gorm"
)

// User 用户模型
// 内嵌 gorm.BaseModel 自动获得 ID、CreatedAt、UpdatedAt、DeletedAt 字段
// gorm tag 用于定义数据库列的类型、约束和索引
type User struct {
	gorm.BaseModel
	Name  string `gorm:"type:varchar(100);not null"`    // 用户名，不允许为空
	Email string `gorm:"type:varchar(100);uniqueIndex"` // 邮箱，唯一索引
	Age   int    `gorm:"default:0"`                     // 年龄，默认值 0
}

func (User) TableName() string {
	return "users"
}

func main() {
	fmt.Println("=== GORM Basic Example ===")
	fmt.Println()

	// 尝试连接 MySQL，使用函数式选项配置连接参数
	// go-boot 的 Option 模式提供灵活且类型安全的参数配置
	fmt.Println("Attempting to connect to MySQL...")
	db, err := gorm.OpenMySQL(
		gorm.WithHost("localhost"),
		gorm.WithPort(3306),
		gorm.WithUser("gate"),
		gorm.WithPassword("123456"),
		gorm.WithDBName("gate"),
		gorm.WithCharset("utf8mb4"),
		gorm.WithParseTime(true),
	)
	if err != nil {
		// 数据库不可用时不阻塞演示，进入 API 展示模式
		fmt.Printf("Cannot connect to MySQL (expected if no DB): %v\n", err)
		fmt.Println("Demo mode: showing GORM API usage without DB connection")
		demoGORMAPI()
		return
	}
	defer db.Close()
	fmt.Println("Connected to MySQL")
}

func demoGORMAPI() {
	fmt.Println("\n=== GORM API Overview ===")
	fmt.Println()
	// 展示支持的数据库类型常量
	fmt.Println("Available DB types:")
	fmt.Printf("  MySQL:      %s\n", gorm.MySQL)
	fmt.Printf("  PostgreSQL: %s\n", gorm.PostgreSQL)
	fmt.Printf("  SQLServer:  %s\n", gorm.SQLServer)
	fmt.Printf("  SQLite:     %s\n", gorm.SQLite)

	// 展示可用的 Option 函数，用于配置数据库连接
	fmt.Println("\nOption functions:")
	fmt.Println("  WithDSN, WithDBType, WithHost, WithPort")
	fmt.Println("  WithUser, WithPassword, WithDBName")
	fmt.Println("  WithMaxIdleConns, WithMaxOpenConns")
	fmt.Println("  WithShowSQL, WithSSLMode, WithCharset, WithTimeZone, WithParseTime")

	// 展示 Open 函数，每种数据库类型对应一个便捷函数
	fmt.Println("\nOpen functions:")
	fmt.Println("  Open(opts...)          - generic")
	fmt.Println("  OpenMySQL(opts...)     - MySQL")
	fmt.Println("  OpenPostgreSQL(opts...) - PostgreSQL")
	fmt.Println("  OpenSQLServer(opts...) - SQL Server")
	fmt.Println("  OpenSQLite(opts...)    - SQLite")

	// 展示 Starter 函数，用于应用启动阶段自动初始化
	fmt.Println("\nStarter functions:")
	fmt.Println("  NewStarter(db, models)           - ping check")
	fmt.Println("  NewAutoMigrateStarter(db, models) - with migration")

	// 展示泛型 Repository[T] 的方法集，提供类型安全的 CRUD
	fmt.Println("\nRepository[T] methods:")
	fmt.Println("  Create, CreateBatch, Update, UpdateByCondition")
	fmt.Println("  Delete, DeleteByCondition")
	fmt.Println("  FindByID, FindOne, FindAll, Count, Raw")

	// 展示 Config 结构的 DSN 生成器，自动构建连接字符串
	fmt.Println("\nConfig DSN generators:")
	cfg := &gorm.Config{
		Host: "localhost", Port: 3306, User: "gate",
		Password: "123456", DBName: "gate", Charset: "utf8mb4",
	}
	fmt.Printf("  MySQL DSN: %s\n", cfg.DSNForMySQL())

	fmt.Println("\ngorm-basic example completed!")
}
