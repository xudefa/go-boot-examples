// Package main 演示 go-boot 与 Casbin + Gin 的访问控制集成
//
// 本示例展示如何:
//   - 创建 Casbin 权限管理器（RBAC 模型）
//   - 在 Gin 路由上应用 Casbin 授权中间件
//   - 根据角色（admin/user）控制资源访问权限
//
// 使用方式:
//
//	cd examples/casbin-gin
//	go run .
//	# Alice 读取 data1 — 应返回 200
//	# curl -H 'X-User: alice' http://localhost:8090/api/data1
//	# Bob 写入 data1 — 应返回 403
//	# curl -H 'X-User: bob' http://localhost:8090/api/data1/write
package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"

	casbin "github.com/xudefa/go-boot-casbin/casbin"
	"github.com/xudefa/go-boot/core"
	ggin "github.com/xudefa/go-boot-gin/server"
	"github.com/xudefa/go-boot/net"
)

func main() {
	fmt.Println("=== Casbin + Gin Example ===")
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

	// 第3步: 创建 Gin 服务器
	server := ggin.New(
		ggin.WithContainer(container),
		ggin.WithMode("debug"),
		ggin.WithPort(8090),
	)
	fmt.Println("Gin server created")

	// 第4步: 创建受 Casbin 保护的路由组，注册自定义授权中间件
	protected := server.Group("/api")
	protected.Use(ggin.AdaptMiddleware(casbin.Authorize(enforcer,
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
		protected.GET("/data1", func(c *gin.Context) {
			c.JSON(200, gin.H{"data": "data1 content", "user": c.GetHeader("X-User")})
		})
		protected.GET("/data1/write", func(c *gin.Context) {
			c.JSON(200, gin.H{"data": "data1 write result", "user": c.GetHeader("X-User")})
		})
		protected.GET("/data2", func(c *gin.Context) {
			c.JSON(200, gin.H{"data": "data2 content", "user": c.GetHeader("X-User")})
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
	fmt.Println("  curl -H 'X-User: alice' http://localhost:8090/api/data1")
	fmt.Println()
	fmt.Println("  # Bob writes data1  — should return 403")
	fmt.Println("  curl -H 'X-User: bob' http://localhost:8090/api/data1/write")
	fmt.Println()

	fmt.Println("Starting server on :8090...")
	fmt.Println("Press Ctrl+C to stop")

	// 第5步: 启动服务器，开始监听端口并处理请求
	if err := server.Start(); err != nil {
		log.Fatalf("服务器错误: %v", err)
	}
}
