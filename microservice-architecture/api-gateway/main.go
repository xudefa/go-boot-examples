package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xudefa/go-boot/core"
	"github.com/xudefa/go-boot/event"
	ginadapter "github.com/xudefa/go-boot-gin/server"
)

// APIService 网关服务，负责转发请求到相应的微服务
type APIService struct {
	httpClient *http.Client
}

func NewAPIService() *APIService {
	return &APIService{
		httpClient: &http.Client{},
	}
}

// ForwardRequest 将请求转发到目标服务
func (s *APIService) ForwardRequest(targetURL string, req *http.Request) (*http.Response, error) {
	// 读取原始请求体
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// 创建新请求
	newReq, err := http.NewRequest(req.Method, targetURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// 复制请求头
	for key, values := range req.Header {
		for _, value := range values {
			newReq.Header.Add(key, value)
		}
	}

	// 发送请求
	resp, err := s.httpClient.Do(newReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func main() {
	// 创建容器
	container := core.New()

	// 创建事件总线
	eventBus := event.NewEventBus()

	// 创建服务
	apiSvc := NewAPIService()

	// 注册服务到容器
	container.Register("apiService", core.Bean(apiSvc))
	container.Register("eventBus", core.Bean(eventBus))

	// 创建Gin实例并绑定容器
	g := ginadapter.New(ginadapter.WithContainer(container))

	// 设置路由
	g.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "api-gateway"})
	})

	// 用户服务路由
	g.Any("/users/*action", func(c *gin.Context) {
		serviceURL := "http://localhost:8083" + strings.TrimPrefix(c.Request.URL.Path, "/users")

		container := g.Container().(core.Container)
		apiSvcRaw, err := container.Get("apiService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		apiSvc := apiSvcRaw.(*APIService)
		resp, err := apiSvc.ForwardRequest(serviceURL, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// 复制响应头
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// 读取并返回响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	})

	// 订单服务路由
	g.Any("/orders/*action", func(c *gin.Context) {
		serviceURL := "http://localhost:8084" + strings.TrimPrefix(c.Request.URL.Path, "/orders")

		container := g.Container().(core.Container)
		apiSvcRaw, err := container.Get("apiService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		apiSvc := apiSvcRaw.(*APIService)
		resp, err := apiSvc.ForwardRequest(serviceURL, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// 复制响应头
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// 读取并返回响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	})

	// 主页路由
	g.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Microservice Architecture API Gateway",
			"services": map[string]string{
				"user-service":  "http://localhost:8083",
				"order-service": "http://localhost:8084",
				"gateway":       "http://localhost:8080",
			},
			"endpoints": map[string]string{
				"users":  "/users/*",
				"orders": "/orders/*",
				"health": "/health",
			},
		})
	})

	fmt.Println("API网关启动在 :8080")

	// 启动服务器
	if err := g.Start(); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}
