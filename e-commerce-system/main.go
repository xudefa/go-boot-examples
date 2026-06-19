package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xudefa/go-boot/aop"
	"github.com/xudefa/go-boot/core"
	"github.com/xudefa/go-boot/event"
	ggin "github.com/xudefa/go-boot-gin/server"
	"github.com/xudefa/go-boot/schedule"
)

// User 用户模型
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"` // admin, customer
}

// Product 商品模型
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

// Cart 购物车项模型
type CartItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// Cart 购物车模型
type Cart struct {
	UserID int        `json:"user_id"`
	Items  []CartItem `json:"items"`
}

// Order 订单模型
type Order struct {
	ID          int         `json:"id"`
	UserID      int         `json:"user_id"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"` // pending, paid, shipped, delivered, cancelled
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// OrderItem 订单项模型
type OrderItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// UserService 用户服务
type UserService struct {
	users  map[int]*User
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (s *UserService) CreateUser(username, email, role string) *User {
	user := &User{
		ID:       s.nextID,
		Username: username,
		Email:    email,
		Role:     role,
	}
	s.users[s.nextID] = user
	s.nextID++
	return user
}

func (s *UserService) GetUser(id int) *User {
	return s.users[id]
}

func (s *UserService) GetUserByUsername(username string) *User {
	for _, user := range s.users {
		if user.Username == username {
			return user
		}
	}
	return nil
}

// ProductService 商品服务
type ProductService struct {
	products map[int]*Product
	nextID   int
}

func NewProductService() *ProductService {
	return &ProductService{
		products: make(map[int]*Product),
		nextID:   1,
	}
}

func (s *ProductService) CreateProduct(name, description string, price float64, stock int, category string) *Product {
	product := &Product{
		ID:          s.nextID,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		Category:    category,
		CreatedAt:   time.Now(),
	}
	s.products[s.nextID] = product
	s.nextID++
	return product
}

func (s *ProductService) GetProduct(id int) *Product {
	return s.products[id]
}

func (s *ProductService) UpdateStock(id int, quantity int) bool {
	product, exists := s.products[id]
	if !exists {
		return false
	}

	newStock := product.Stock + quantity
	if newStock < 0 {
		return false // 库存不足
	}

	product.Stock = newStock
	return true
}

func (s *ProductService) GetAllProducts() []*Product {
	products := make([]*Product, 0, len(s.products))
	for _, product := range s.products {
		products = append(products, product)
	}
	return products
}

// CartService 购物车服务
type CartService struct {
	carts map[int]*Cart
}

func NewCartService() *CartService {
	return &CartService{
		carts: make(map[int]*Cart),
	}
}

func (s *CartService) AddToCart(userID, productID, quantity int) error {
	cart, exists := s.cart(userID)
	if !exists {
		cart = &Cart{UserID: userID, Items: []CartItem{}}
		s.carts[userID] = cart
	}

	// 检查是否已存在该商品
	itemExists := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].Quantity += quantity
			itemExists = true
			break
		}
	}

	if !itemExists {
		cart.Items = append(cart.Items, CartItem{
			ProductID: productID,
			Quantity:  quantity,
		})
	}

	return nil
}

func (s *CartService) GetCart(userID int) *Cart {
	cart, exists := s.carts[userID]
	if !exists {
		return &Cart{UserID: userID, Items: []CartItem{}}
	}
	return cart
}

func (s *CartService) RemoveFromCart(userID, productID int) error {
	cart, exists := s.carts[userID]
	if !exists {
		return fmt.Errorf("cart not found")
	}

	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("product not found in cart")
}

func (s *CartService) ClearCart(userID int) error {
	delete(s.carts, userID)
	return nil
}

func (s *CartService) cart(userID int) (*Cart, bool) {
	cart, exists := s.carts[userID]
	return cart, exists
}

// OrderService 订单服务
type OrderService struct {
	orders     map[int]*Order
	nextID     int
	userSvc    *UserService
	productSvc *ProductService
}

func NewOrderService(userSvc *UserService, productSvc *ProductService) *OrderService {
	return &OrderService{
		orders:     make(map[int]*Order),
		nextID:     1,
		userSvc:    userSvc,
		productSvc: productSvc,
	}
}

func (s *OrderService) CreateOrder(userID int, items []OrderItem) (*Order, error) {
	user := s.userSvc.GetUser(userID)
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	var totalAmount float64
	orderItems := make([]OrderItem, len(items))

	for i, item := range items {
		product := s.productSvc.GetProduct(item.ProductID)
		if product == nil {
			return nil, fmt.Errorf("product %d not found", item.ProductID)
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %d", item.ProductID)
		}

		orderItems[i] = OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		totalAmount += product.Price * float64(item.Quantity)
	}

	order := &Order{
		ID:          s.nextID,
		UserID:      userID,
		Items:       orderItems,
		TotalAmount: math.Round(totalAmount*100) / 100, // 保留两位小数
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.orders[s.nextID] = order
	s.nextID++

	// 扣减库存
	for _, item := range items {
		s.productSvc.UpdateStock(item.ProductID, -item.Quantity)
	}

	return order, nil
}

func (s *OrderService) GetOrder(id int) *Order {
	return s.orders[id]
}

func (s *OrderService) UpdateOrderStatus(orderID int, status string) error {
	order, exists := s.orders[orderID]
	if !exists {
		return fmt.Errorf("order not found")
	}

	order.Status = status
	order.UpdatedAt = time.Now()
	return nil
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

// InventoryMonitorService 库存监控服务
type InventoryMonitorService struct {
	productSvc *ProductService
	eventBus   *event.EventBus
}

func NewInventoryMonitorService(productSvc *ProductService, eventBus *event.EventBus) *InventoryMonitorService {
	return &InventoryMonitorService{
		productSvc: productSvc,
		eventBus:   eventBus,
	}
}

// CheckLowStock 检查低库存商品
func (s *InventoryMonitorService) CheckLowStock(threshold int) []*Product {
	var lowStockProducts []*Product
	products := s.productSvc.GetAllProducts()

	for _, product := range products {
		if product.Stock <= threshold {
			lowStockProducts = append(lowStockProducts, product)
		}
	}

	return lowStockProducts
}

// InventoryAlertEvent 库存警报事件
type InventoryAlertEvent struct {
	EventType    string
	ProductID    int
	ProductName  string
	CurrentStock int
	EventTime    time.Time
}

func (e *InventoryAlertEvent) Type() string {
	return e.EventType
}

func (e *InventoryAlertEvent) Timestamp() time.Time {
	if e.EventTime.IsZero() {
		return time.Now()
	}
	return e.EventTime
}

// AOP 日志切面
type LoggingAspect struct{}

func (l *LoggingAspect) Before(jp aop.JoinPoint) {
	fmt.Printf("[LOG] Before method: %s, Args: %v\n", jp.Signature().Name(), jp.Args())
}

func (l *LoggingAspect) AfterReturning(jp aop.JoinPoint, result any) {
	fmt.Printf("[LOG] After method: %s, Result: %v\n", jp.Signature().Name(), result)
}

func (l *LoggingAspect) AfterThrowing(jp aop.JoinPoint, err error) {
	fmt.Printf("[LOG] Method: %s, Error: %v\n", jp.Signature().Name(), err)
}

func main() {
	// 创建容器
	container := core.New()

	// 创建事件总线
	eventBus := event.NewEventBus()

	// 创建服务
	userSvc := NewUserService()
	productSvc := NewProductService()
	cartSvc := NewCartService()
	orderSvc := NewOrderService(userSvc, productSvc)
	inventoryMonitorSvc := NewInventoryMonitorService(productSvc, eventBus)

	// 注册服务到容器
	container.Register("userService", core.Bean(userSvc))
	container.Register("productService", core.Bean(productSvc))
	container.Register("cartService", core.Bean(cartSvc))
	container.Register("orderService", core.Bean(orderSvc))
	container.Register("inventoryMonitorService", core.Bean(inventoryMonitorSvc))
	container.Register("eventBus", core.Bean(eventBus))

	// 创建AOP代理
	productProxyFactory := aop.NewProxyFactory(productSvc)
	productProxyFactory.SetAspects([]*aop.AspectMeta{
		{
			PointCut: aop.MatchByName("CreateProduct"),
			Advice:   aop.Around(new(LoggingAspect).Around),
		},
		{
			PointCut: aop.MatchByName("UpdateStock"),
			Advice:   aop.Around(new(LoggingAspect).Around),
		},
	})
	productSvcProxy := productProxyFactory.GetProxy()
	container.Register("productServiceProxy", core.Bean(productSvcProxy))

	// 创建Gin实例并绑定容器
	g := ggin.New(ggin.WithContainer(container))

	// 设置路由
	setupRoutes(g)

	// 启动定时任务
	go startScheduler(inventoryMonitorSvc, eventBus)

	fmt.Println("电商系统启动在 :8080")
	if err := g.Start(); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func (l *LoggingAspect) Around(jp aop.JoinPoint, proceed aop.ProceedFunc) any {
	l.Before(jp)

	result := proceed(jp.Args())

	// 检查是否有错误返回
	if result != nil {
		if err, ok := result.(error); ok {
			l.AfterThrowing(jp, err)
			return result
		}
	}

	l.AfterReturning(jp, result)
	return result
}

func setupRoutes(g *ggin.GinServer) {
	// 用户相关路由
	g.POST("/api/users", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Role == "" {
			req.Role = "customer"
		}

		container := g.Container().(core.Container)
		userSvcRaw, err := container.Get("userService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userSvc := userSvcRaw.(*UserService)
		user := userSvc.CreateUser(req.Username, req.Email, req.Role)

		c.JSON(http.StatusOK, user)
	})

	g.GET("/api/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		container := g.Container().(core.Container)
		userSvcRaw, err := container.Get("userService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userSvc := userSvcRaw.(*UserService)
		user := userSvc.GetUser(id)
		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	// 商品相关路由
	g.GET("/api/products", func(c *gin.Context) {
		container := g.Container().(core.Container)
		productSvcRaw, err := container.Get("productService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		productSvc := productSvcRaw.(*ProductService)
		products := productSvc.GetAllProducts()

		c.JSON(http.StatusOK, products)
	})

	g.GET("/api/products/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		container := g.Container().(core.Container)
		productSvcRaw, err := container.Get("productService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		productSvc := productSvcRaw.(*ProductService)
		product := productSvc.GetProduct(id)
		if product == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		c.JSON(http.StatusOK, product)
	})

	g.POST("/api/products", func(c *gin.Context) {
		var req struct {
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Price       float64 `json:"price"`
			Stock       int     `json:"stock"`
			Category    string  `json:"category"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		container := g.Container().(core.Container)
		productSvcRaw, err := container.Get("productService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		productSvc := productSvcRaw.(*ProductService)
		product := productSvc.CreateProduct(req.Name, req.Description, req.Price, req.Stock, req.Category)

		c.JSON(http.StatusOK, product)
	})

	// 购物车相关路由
	g.GET("/api/cart/:userId", func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("userId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		container := g.Container().(core.Container)
		cartSvcRaw, err := container.Get("cartService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cartSvc := cartSvcRaw.(*CartService)
		cart := cartSvc.GetCart(userId)

		c.JSON(http.StatusOK, cart)
	})

	g.POST("/api/cart/add", func(c *gin.Context) {
		var req struct {
			UserID    int `json:"user_id"`
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		container := g.Container().(core.Container)
		cartSvcRaw, err := container.Get("cartService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cartSvc := cartSvcRaw.(*CartService)
		err = cartSvc.AddToCart(req.UserID, req.ProductID, req.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Added to cart"})
	})

	g.DELETE("/api/cart/remove", func(c *gin.Context) {
		var req struct {
			UserID    int `json:"user_id"`
			ProductID int `json:"product_id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		container := g.Container().(core.Container)
		cartSvcRaw, err := container.Get("cartService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cartSvc := cartSvcRaw.(*CartService)
		removeErr := cartSvc.RemoveFromCart(req.UserID, req.ProductID)
		if removeErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": removeErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Removed from cart"})
	})

	// 订单相关路由
	g.POST("/api/orders", func(c *gin.Context) {
		var req struct {
			UserID int         `json:"user_id"`
			Items  []OrderItem `json:"items"`
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
		order, err := orderSvc.CreateOrder(req.UserID, req.Items)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, order)
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
		err = orderSvc.UpdateOrderStatus(id, req.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
	})

	// 库存监控相关路由
	g.GET("/api/inventory/low-stock", func(c *gin.Context) {
		container := g.Container().(core.Container)
		inventorySvcRaw, err := container.Get("inventoryMonitorService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		inventorySvc := inventorySvcRaw.(*InventoryMonitorService)
		var threshold int
		if thresholdStr := c.Query("threshold"); thresholdStr != "" {
			if t, err := strconv.Atoi(thresholdStr); err == nil {
				threshold = t
			} else {
				threshold = 10 // 默认阈值
			}
		} else {
			threshold = 10 // 默认阈值
		}

		products := inventorySvc.CheckLowStock(threshold)
		c.JSON(http.StatusOK, products)
	})
}

func startScheduler(inventorySvc *InventoryMonitorService, eventBus *event.EventBus) {
	scheduler := schedule.NewScheduler()

	// 每分钟检查一次低库存
	ctx := context.Background()
	task1 := schedule.NewTask("check-low-stock", "* * * * * ?", func(ctx context.Context) error {
		fmt.Println("Checking low stock...")
		lowStockProducts := inventorySvc.CheckLowStock(5)

		for _, product := range lowStockProducts {
			alertEvent := &InventoryAlertEvent{
				EventType:    "inventory.alert",
				ProductID:    product.ID,
				ProductName:  product.Name,
				CurrentStock: product.Stock,
				EventTime:    time.Now(),
			}
			eventBus.Publish(alertEvent)
			fmt.Printf("Low stock alert: %s (ID: %d, Stock: %d)\n",
				product.Name, product.ID, product.Stock)
		}
		return nil
	})

	scheduler.Register(task1)

	scheduler.Start(ctx)
	// 注意：在实际应用中，这里不应该阻塞等待，但为了示例，我们暂时这样处理
	// 通常我们会使用一个channel来等待退出信号
}
