// Package main 演示 go-boot 与 Casbin + Hertz 的访问控制集成
//
// 本示例展示如何:
//   - 创建 Casbin 权限管理器（RBAC 模型）
//   - 在 Hertz 路由上应用 Casbin 授权中间件
//   - 根据角色（admin/user）控制资源访问权限
//
// 使用方式:
//
//	cd examples/casbin-hertz
//	go run .
//	# Alice 读取 data1 — 应返回 200
//	# curl -H 'X-User: alice' http://localhost:8091/api/data1
//	# Bob 写入 data1 — 应返回 403
//	# curl -H 'X-User: bob' http://localhost:8091/api/data1/write
package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	casbin "github.com/xudefa/go-boot-casbin/casbin"
	hertz "github.com/xudefa/go-boot-hertz/server"
	"github.com/xudefa/go-boot/core"
	"github.com/xudefa/go-boot/net"
)

func main() {
	fmt.Println("=== Casbin + Hertz Example ===")
	fmt.Println()

	// 获取当前文件所在目录
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	modelPath := filepath.Join(dir, "model.conf")
	policyPath := filepath.Join(dir, "policy.csv")
	// 第1步: 创建 Casbin 权限执行器，从本地文件加载 RBAC 模型配置和策略规则
	enforcer, err := casbin.NewEnforcer(
		casbin.WithModel(modelPath),
		casbin.WithAdapter(policyPath),
	)
	if err != nil {
		log.Fatalf("创建执行器失败: %v", err)
	}
	fmt.Println("Casbin enforcer created with RBAC model")

	// 第2步: 创建核心容器并注册 Casbin 执行器
	container := core.New()
	container.Register(casbin.EnforcerBeanID, core.Bean(enforcer))
	fmt.Println("Enforcer registered in container")

	// 第3步: 创建 Hertz 服务器
	server := hertz.NewServer(
		hertz.WithServerContainer(container),
		hertz.WithPort(8091),
	)
	fmt.Println("Hertz server created")

	// 第4步: 创建受 Casbin 保护的路由组，注册自定义授权中间件
	engine := server.Engine()
	api := engine.Group("/api")
	api.Use(server.AdaptMiddleware(casbin.Authorize(enforcer,
		func(ctx net.HandlerContext) (string, string, string) {
			sub := ctx.Header("X-User")
			obj := strings.TrimPrefix(ctx.RequestURI(), "/api/")
			act := "read"
			if strings.HasSuffix(obj, "/write") {
				act = "write"
				obj = strings.TrimSuffix(obj, "/write")
			}
			return sub, obj, act
		},
	)))
	{
		api.GET("/data1", func(ctx context.Context, c *app.RequestContext) {
			c.JSON(200, map[string]interface{}{
				"data": "data1 content",
				"user": string(c.GetHeader("X-User")),
			})
		})
		api.GET("/data1/write", func(ctx context.Context, c *app.RequestContext) {
			c.JSON(200, map[string]interface{}{
				"data": "data1 write result",
				"user": string(c.GetHeader("X-User")),
			})
		})
		api.GET("/data2", func(ctx context.Context, c *app.RequestContext) {
			c.JSON(200, map[string]interface{}{
				"data": "data2 content",
				"user": string(c.GetHeader("X-User")),
			})
		})
	}
	fmt.Println("Routes registered with Casbin authorization:")
	fmt.Println("  GET /api/data1       - requires read access")
	fmt.Println("  GET /api/data1/write - requires write access")
	fmt.Println("  GET /api/data2       - requires read access")
	fmt.Println()
	fmt.Println("Policy rules:")
	fmt.Println("  admin (alice): read/write data1, data2")
	fmt.Println("  user  (bob):   read data1, data2")
	fmt.Println()
	fmt.Println("Test commands:")
	fmt.Println("  # Alice reads data1 — should return 200")
	fmt.Println("  curl -H 'X-User: alice' http://localhost:8091/api/data1")
	fmt.Println()
	fmt.Println("  # Bob writes data1  — should return 403")
	fmt.Println("  curl -H 'X-User: bob' http://localhost:8091/api/data1/write")
	fmt.Println()

	fmt.Println("Starting server on :8091...")
	fmt.Println("Press Ctrl+C to stop")

	// 第5步: 启动服务器，开始监听端口并处理请求
	if err := server.Start(); err != nil {
		log.Fatalf("服务器错误: %v", err)
	}
}
