package main

import (
	"fmt"
	"log"

	_ "github.com/xudefa/go-boot-consul"
	_ "github.com/xudefa/go-boot-etcd"
	_ "github.com/xudefa/go-boot-nacos"
	"github.com/xudefa/go-boot/boot"
)

func main() {
	fmt.Println("=== Config Center Example ===")
	fmt.Println("Note: This example requires Nacos/Consul/Etcd running locally")
	fmt.Println()
	// 启用Consul配置中心
	//	boot.WithConfigCenter("consul", []string{"127.0.0.1:8500"},
	// 		boot.WithConfigCenterTimeout(5*time.Second),
	// 	),
	// 或者nacos:
	// boot.WithConfigCenter("nacos", []string{"127.0.0.1:8848"},
	// 	boot.WithConfigCenterDataID("com.xdf.kafka.SendMessageService:1.0.0::provider:kafka"),
	// 	boot.WithConfigCenterGroup("dubbo"),
	// 	boot.WithConfigCenterTimeout(5*time.Second),
	// ),

	// 尝试创建带配置中心的应用
	app, err := boot.NewApplication(
		boot.WithAppName("config-center-example"),
		boot.WithVersion("1.0.0"),
		boot.WithProfiles("dev"),
	)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	if err := app.Start(); err != nil {
		fmt.Printf("Warning: Failed to start application: %v\n", err)
		fmt.Println("Running in demo mode without config center...")
		fmt.Println("\nconfig-center example completed!")
		return
	}
	defer app.Stop()

	fmt.Println("Application started with config center support")

	// 修复后，配置会使用从 DataID 提取的有意义键名
	serviceName := app.Environment().GetString("canonicalName", "")
	if serviceName != "" {
		fmt.Println("Kafka SendMessageService:", serviceName)
	} else {
		fmt.Println("Kafka SendMessageService: not found in config (expected without running config center)")
	}

	fmt.Println("\nconfig-center example completed!")
}
