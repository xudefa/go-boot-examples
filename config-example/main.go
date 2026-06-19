// Package main 演示 go-boot 的配置管理抽象
//
// 本示例展示如何:
//   - 使用 ConfigModel 与 ConfigOption
//   - 使用 Validator 进行配置校验
//   - 使用 WatchManager 监听配置变更
//
// 使用方式:
//
//	cd examples/config-example
//	go run .
package main

import (
	"fmt"

	"github.com/xudefa/go-boot/config"
)

func main() {
	fmt.Println("=== Config Abstraction Example ===")
	fmt.Println()

	// 打印 Config 接口定义的所有方法，展示配置抽象层的完整 API 表面
	fmt.Println("=== Config interface ===")
	fmt.Println("The Config interface defines:")
	fmt.Println("  Get, GetAll, GetString, GetInt, GetBool, GetFloat64")
	fmt.Println("  GetStringMap, GetStringSlice, GetIntSlice")
	fmt.Println("  HasKey, Unmarshal, UnmarshalKey")
	fmt.Println("  Watch, StopWatch, GetSource")

	// 创建 ConfigOption 列表，演示函数式选项模式的用法
	// 可以配置文件名、搜索路径、文件类型、环境名称
	fmt.Println("\n=== Config options ===")
	opts := []config.ConfigOption{
		config.WithConfigName("application"),
		config.WithConfigPath(".", "./config"),
		config.WithConfigType("yaml"),
		config.WithEnvironment("dev"),
	}
	fmt.Printf("Config options created: %d options\n", len(opts))

	// 使用 New 创建 ConfigModel 实例，传入配置选项
	// 第一个参数是 load 函数，这里传 nil 表示不加载
	model, err := config.New(nil, opts...)
	if err != nil {
		fmt.Printf("Config.New: %v\n", err)
	} else {
		fmt.Printf("ConfigModel created: name=%s, type=%s, env=%s\n",
			model.ConfigName, model.ConfigType, model.Env)
	}

	// Validator 支持多种校验规则：
	//   AddRequired - 必填字段
	//   AddMin/AddMax - 数值范围
	//   AddRegex - 正则匹配
	//   AddEnum - 枚举允许值
	fmt.Println("\n=== Validator ===")
	validator := config.NewValidator()
	validator.AddRequired("host", "port")
	validator.AddMin("port", 1024)
	validator.AddMax("port", 65535)
	validator.AddRegex("host", `^[a-zA-Z0-9.-]+$`)
	validator.AddEnum("mode", "debug", "release", "test")
	fmt.Println("Validator configured with rules")

	// 验证配置对象，输出通过或失败信息
	// Validator 只接受 map[string]any 类型
	cfg := map[string]any{"host": "localhost", "port": 8080}
	err = validator.Validate(cfg)
	if err != nil {
		fmt.Printf("Validation errors: %v\n", err)
	} else {
		fmt.Printf("Validation passed for %+v\n", cfg)
	}

	// WatchManager 支持注册配置变更回调，当配置修改/删除/创建时通知监听器
	fmt.Println("\n=== WatchManager ===")
	wm := config.NewWatchManager()
	wm.Register("server.port", func(evt config.WatchEvent) {
		fmt.Printf("Watch event: type=%s, key=%s, value=%v\n",
			evt.Type, evt.Key, evt.Value)
	})
	// 模拟修改事件：server.port 变为 9090
	wm.Notify(config.WatchEvent{
		Type:  config.EventModify,
		Key:   "server.port",
		Value: 9090,
	})
	wm.Close()
	fmt.Println("Watch events: modify, delete, create")

	fmt.Println("\nconfig-example completed successfully!")
}
