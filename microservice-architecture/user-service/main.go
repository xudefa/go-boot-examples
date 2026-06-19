package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xudefa/go-boot/core"
	"github.com/xudefa/go-boot/event"
	ginadapter "github.com/xudefa/go-boot-gin/server"
)

// User 用户模型
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Created  string `json:"created"`
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

func (s *UserService) CreateUser(username, email, phone, address string) *User {
	user := &User{
		ID:       s.nextID,
		Username: username,
		Email:    email,
		Phone:    phone,
		Address:  address,
		Created:  "2023-01-01",
	}
	s.users[s.nextID] = user
	s.nextID++
	return user
}

func (s *UserService) GetUser(id int) *User {
	return s.users[id]
}

func (s *UserService) UpdateUser(id int, username, email, phone, address string) *User {
	user, exists := s.users[id]
	if !exists {
		return nil
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	if phone != "" {
		user.Phone = phone
	}
	if address != "" {
		user.Address = address
	}

	return user
}

func (s *UserService) GetAllUsers() []*User {
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func main() {
	// 创建容器
	container := core.New()

	// 创建事件总线
	eventBus := event.NewEventBus()

	// 创建服务
	userSvc := NewUserService()

	// 注册服务到容器
	container.Register("userService", core.Bean(userSvc))
	container.Register("eventBus", core.Bean(eventBus))

	// 创建Gin实例并绑定容器
	g := ginadapter.New(ginadapter.WithContainer(container))

	// 设置路由
	g.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "user-service"})
	})

	g.GET("/api/users", func(c *gin.Context) {
		container := g.Container().(core.Container)
		userSvcRaw, err := container.Get("userService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userSvc := userSvcRaw.(*UserService)
		users := userSvc.GetAllUsers()

		c.JSON(http.StatusOK, users)
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

	g.POST("/api/users", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Address  string `json:"address"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		container := g.Container().(core.Container)
		userSvcRaw, err := container.Get("userService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userSvc := userSvcRaw.(*UserService)
		user := userSvc.CreateUser(req.Username, req.Email, req.Phone, req.Address)

		c.JSON(http.StatusOK, user)
	})

	g.PUT("/api/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Address  string `json:"address"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		container := g.Container().(core.Container)
		userSvcRaw, err := container.Get("userService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userSvc := userSvcRaw.(*UserService)
		user := userSvc.UpdateUser(id, req.Username, req.Email, req.Phone, req.Address)
		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	fmt.Println("用户服务启动在 :8083")

	// 启动服务器
	if err := g.Start(); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}
