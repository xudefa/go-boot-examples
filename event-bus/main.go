// Package main 演示 go-boot 的事件驱动支持
//
// 本示例展示如何:
//   - 创建事件总线
//   - 订阅和发布事件
//   - 使用内置事件类型
//
// 使用方式:
//
//	cd examples/event-bus
//	go run .
package main

import (
	"fmt"
	"time"

	"github.com/xudefa/go-boot/event"
)

// UserRegisteredEvent 自定义用户注册事件
type UserRegisteredEvent struct {
	event.BaseEvent
	UserID   int
	UserName string
}

func main() {
	fmt.Println("=== Event Bus Example ===")
	fmt.Println()

	// NewEventBus 创建事件总线实例，管理发布/订阅
	bus := event.NewEventBus()
	fmt.Println("Event bus created")

	// 订阅 UserRegistered 事件，注册处理器 1：打印用户注册信息
	bus.Subscribe("UserRegistered", func(e event.ApplicationEvent) {
		ue, ok := e.(*UserRegisteredEvent)
		if !ok {
			fmt.Println("Received event of wrong type")
			return
		}
		fmt.Printf("  [Handler 1] User registered: ID=%d, Name=%s\n", ue.UserID, ue.UserName)
	})

	// 订阅 UserRegistered 事件，注册处理器 2：发送欢迎邮件
	// 多个订阅者独立运行，互不干扰
	bus.Subscribe("UserRegistered", func(e event.ApplicationEvent) {
		ue, _ := e.(*UserRegisteredEvent)
		fmt.Printf("  [Handler 2] Send welcome email to user %d\n", ue.UserID)
	})

	fmt.Println("\nPublishing UserRegistered event:")
	// 创建自定义事件实例，设置事件类型、时间戳和业务数据
	evt := &UserRegisteredEvent{
		BaseEvent: event.BaseEvent{
			EventType: "UserRegistered",
			EventTime: time.Now(),
		},
		UserID:   1001,
		UserName: "Alice",
	}
	// Publish 广播事件到所有订阅者，每个处理器都会收到通知
	bus.Publish(evt)

	// 内置事件类型对应应用启动的不同阶段，贯穿整个应用生命周期
	fmt.Println("\n=== Built-in event types ===")
	fmt.Printf("  EventEnvironmentPrepared: %s\n", event.EventEnvironmentPrepared)
	fmt.Printf("  EventContextRefreshed: %s\n", event.EventContextRefreshed)
	fmt.Printf("  EventApplicationStarted: %s\n", event.EventApplicationStarted)
	fmt.Printf("  EventApplicationReady: %s\n", event.EventApplicationReady)
	fmt.Printf("  EventApplicationStopped: %s\n", event.EventApplicationStopped)

	fmt.Println("\nevent-bus example completed successfully!")
}
