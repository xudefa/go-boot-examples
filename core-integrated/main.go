// Package main 演示 go-boot 核心 - 集成使用与依赖
//
// 本示例展示如何:
//   - 注册多个带依赖的 Bean
//   - 使用 core.DependsOn() 显式声明依赖
//   - 使用 core.Init() 进行初始化
//   - 使用 core.Condition() 条件注册
//
// 使用方式:
//
//	cd examples/core-integrated
//	go run .
package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/xudefa/go-boot/core"
)

// Logger 接口
type Logger interface {
	Info(msg string)
	Error(msg string)
}

// SimpleLogger 实现
type SimpleLogger struct {
	Name string
}

func (l *SimpleLogger) Info(msg string) {
	fmt.Printf("[INFO] %s: %s\n", l.Name, msg)
}

func (l *SimpleLogger) Error(msg string) {
	fmt.Printf("[ERROR] %s: %s\n", l.Name, msg)
}

// Database 结构体
type Database struct {
	URL       string
	Connected bool
}

// UserService 带依赖
type UserService struct {
	// 注入的字段
	Log *SimpleLogger `inject:"logger"`
	DB  *Database     `inject:"database"`

	Name string
}

func (s *UserService) Init() error {
	fmt.Println("  UserService.Init() called")
	s.Name = "InitializedService"
	return nil
}

func main() {
	// 第1步: 创建容器
	container := core.New()
	fmt.Println("=== Core Integrated Example ===")
	fmt.Println()

	// 第2步: 注册日志器
	// 日志器是最基础的依赖，后续的 Bean 可能会引用它
	err := container.Register("logger",
		core.Bean(&SimpleLogger{Name: "AppLogger"}),
		core.Singleton(),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Registered: logger")

	// 第3步: 使用 Condition 注册数据库
	// Condition 确保只有 logger 存在时才会创建数据库 Bean
	err = container.Register("database",
		core.Factory(
			func(c core.Container) (any, error) {
				fmt.Println("  Creating database connection...")
				time.Sleep(100 * time.Millisecond)
				return &Database{URL: "localhost:5432", Connected: true}, nil
			},
			reflect.TypeFor[Database](),
		),
		core.Singleton(),
		core.Condition(func(c core.Container) bool {
			return c.Has("logger")
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Registered: database (with condition)")

	// 第4步: 注册带依赖和初始化的用户服务
	err = container.Register("userService",
		core.Bean(&UserService{}),
		core.DependsOn("logger", "database"),
		core.Init(func(bean any) error {
			return bean.(*UserService).Init()
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Registered: userService (with dependencies and init)")
	fmt.Println()

	// 第5步: 获取并使用服务
	// DependsOn 确保 logger 和 database 先于 userService 初始化
	bean, err := container.Get("userService")
	if err != nil {
		log.Fatal(err)
	}

	service := bean.(*UserService)
	fmt.Println("Retrieved userService:")
	fmt.Printf("  Name: %s\n", service.Name)
	fmt.Printf("  Logger: %v (injected)\n", service.Log != nil)
	fmt.Printf("  Database: %v (injected, connected=%v)\n", service.DB != nil, service.DB.Connected)
	fmt.Println()

	// 第6步: 使用服务的依赖
	service.Log.Info("User service is working!")
	service.Log.Error("This is an error message")

	// 第7步: 演示 Invoke
	// Invoke 自动从容器中解析参数并调用函数，无需手动 Get
	fmt.Println("Using container.Invoke():")
	_, err = container.Invoke(func(log *SimpleLogger, db *Database) any {
		log.Info(fmt.Sprintf("Invoked with db URL: %s", db.URL))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\ncore-integrated example completed successfully!")
}
