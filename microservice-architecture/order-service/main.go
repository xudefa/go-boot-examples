package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xudefa/go-boot/core"
	"github.com/xudefa/go-boot/event"
	ginadapter "github.com/xudefa/go-boot-gin/server"
)

// Order 订单模型
type Order struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ProductIDs  []int     `json:"product_ids"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"` // pending, paid, shipped, delivered, cancelled
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrderService 订单服务
type OrderService struct {
	orders     map[int]*Order
	nextID     int
	httpClient *http.Client
}

func NewOrderService() *OrderService {
	return &OrderService{
		orders:     make(map[int]*Order),
		nextID:     1,
		httpClient: &http.Client{},
	}
}

func (s *OrderService) CreateOrder(userID int, productIDs []int, totalAmount float64) *Order {
	order := &Order{
		ID:          s.nextID,
		UserID:      userID,
		ProductIDs:  productIDs,
		TotalAmount: totalAmount,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.orders[s.nextID] = order
	s.nextID++
	return order
}

func (s *OrderService) GetOrder(id int) *Order {
	return s.orders[id]
}

func (s *OrderService) UpdateOrderStatus(id int, status string) *Order {
	order, exists := s.orders[id]
	if !exists {
		return nil
	}

	order.Status = status
	order.UpdatedAt = time.Now()
	return order
}

func (s *OrderService) GetUserOrders(userID int) []*Order {
	var orders []*Order
	for _, order := range s.orders {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}
	return orders
}

// GetUser 从用户服务获取用户信息
func (s *OrderService) GetUser(userID int) (*User, error) {
	resp, err := s.httpClient.Get(fmt.Sprintf("http://localhost:8083/api/users/%d", userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: %s", string(body))
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// User 用户模型（从用户服务获取）
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Created  string `json:"created"`
}

func main() {
	// 创建容器
	container := core.New()

	// 创建事件总线
	eventBus := event.NewEventBus()

	// 创建服务
	orderSvc := NewOrderService()

	// 注册服务到容器
	container.Register("orderService", core.Bean(orderSvc))
	container.Register("eventBus", core.Bean(eventBus))

	// 创建Gin实例并绑定容器
	g := ginadapter.New(ginadapter.WithContainer(container))

	// 设置路由
	g.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "order-service"})
	})

	g.GET("/api/orders", func(c *gin.Context) {
		container := g.Container().(core.Container)
		orderSvcRaw, err := container.Get("orderService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orderSvc := orderSvcRaw.(*OrderService)
		// 返回所有订单
		orders := make([]*Order, 0, len(orderSvc.orders))
		for _, order := range orderSvc.orders {
			orders = append(orders, order)
		}

		c.JSON(http.StatusOK, orders)
	})

	g.GET("/api/orders/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		container := g.Container().(core.Container)
		orderSvcRaw, err := container.Get("orderService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orderSvc := orderSvcRaw.(*OrderService)
		order := orderSvc.GetOrder(id)
		if order == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, order)
	})

	g.POST("/api/orders", func(c *gin.Context) {
		var req struct {
			UserID      int     `json:"user_id"`
			ProductIDs  []int   `json:"product_ids"`
			TotalAmount float64 `json:"total_amount"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证用户是否存在
		container := g.Container().(core.Container)
		orderSvcRaw, err := container.Get("orderService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orderSvc := orderSvcRaw.(*OrderService)
		_, err = orderSvc.GetUser(req.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("User validation failed: %v", err)})
			return
		}

		order := orderSvc.CreateOrder(req.UserID, req.ProductIDs, req.TotalAmount)

		c.JSON(http.StatusOK, order)
	})

	g.PUT("/api/orders/:id/status", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		var req struct {
			Status string `json:"status"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		container := g.Container().(core.Container)
		orderSvcRaw, err := container.Get("orderService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orderSvc := orderSvcRaw.(*OrderService)
		order := orderSvc.UpdateOrderStatus(id, req.Status)
		if order == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, order)
	})

	g.GET("/api/users/:id/orders", func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		container := g.Container().(core.Container)
		orderSvcRaw, err := container.Get("orderService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orderSvc := orderSvcRaw.(*OrderService)
		orders := orderSvc.GetUserOrders(userId)

		c.JSON(http.StatusOK, orders)
	})

	fmt.Println("订单服务启动在 :8084")

	// 启动服务器
	if err := g.Start(); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}
