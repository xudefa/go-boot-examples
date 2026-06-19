// Package main 演示 go-boot 与 Viper - 配置管理
//
// 本示例展示如何:
//   - 使用 ViperConfig 创建配置管理器
//   - 使用 FileLoader 加载配置文件
//   - 使用默认值和环境变量
//
// 使用方式:
//
//	cd examples/viper-config
//	go run .
package main

import (
	"context"
	"fmt"
	"os"

	viper "github.com/xudefa/go-boot-viper"
	"github.com/xudefa/go-boot/config"
)

func main() {
	fmt.Println("=== Viper Config Example ===")
	fmt.Println()

	// 创建 ViperConfig 实例，指定配置文件名、路径和类型
	// 当配置文件不存在时会返回错误，但此时可用默认值继续运行
	cfg, err := viper.New(
		config.WithConfigName("config"),
		config.WithConfigPath("."),
		config.WithConfigType("yaml"),
	)
	if err != nil {
		fmt.Printf("Note: No config file found, using defaults: %v\n", err)
	} else {
		// 通过 GetString 读取 server.host 和 server.mode 配置项
		host := cfg.GetString("server.host")
		mode := cfg.GetString("server.mode")
		fmt.Printf("Loaded config - host: %s, mode: %s\n", host, mode)
	}

	// FileLoader 支持环境感知的配置文件加载
	// 可通过 WithLoaderEnv 传入环境变量，加载对应的 profile 配置文件
	fmt.Println("\n=== Using file loader ===")
	loader := viper.NewFileLoader()
	cfgi, err := loader.Load(
		config.WithFileName("config"),
		config.WithLoaderEnv(os.Getenv("APP_ENV")),
	)
	if err != nil {
		fmt.Printf("FileLoader: config file not found (expected): %v\n", err)
	} else {
		// 类型断言获取 ViperConfig，调用 GetAll 查看所有配置
		vc := cfgi.(*viper.ViperConfig)
		fmt.Printf("Config loaded: %+v\n", vc.GetAll())
	}

	// MustNew 不返回错误，适合在初始化时使用默认配置启动
	// SetDefault 设置单个键的默认值，SetDefaults 批量设置
	fmt.Println("\n=== MustNew with defaults ===")
	cfg2 := viper.MustNew()
	cfg2.SetDefault("server.port", 8080)
	cfg2.SetDefaults(map[string]any{"app.name": "go-boot-demo"})
	fmt.Printf("Default port: %d\n", cfg2.GetInt("server.port"))
	fmt.Printf("Default app name: %s\n", cfg2.GetString("app.name"))

	// NewWithContext 支持传入 context.Context，便于链路追踪和超时控制
	fmt.Println("\n=== Context-based creation ===")
	cfg3, err := viper.NewWithContext(context.Background())
	if err != nil {
		fmt.Printf("Context-based creation: %v\n", err)
	} else {
		fmt.Printf("Config created with context: %s\n", cfg3.GetSource())
	}

	fmt.Println("\nviper-config example completed successfully!")
}
