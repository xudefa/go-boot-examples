package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xudefa/go-boot/core"
	"github.com/xudefa/go-boot/event"
	ginadapter "github.com/xudefa/go-boot-gin/server"
)

// User 用户模型
type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Online   bool      `json:"online"`
	JoinedAt time.Time `json:"joined_at"`
}

// Message 消息模型
type Message struct {
	ID          int       `json:"id"`
	SenderID    int       `json:"sender_id"`
	Sender      string    `json:"sender"`
	Room        string    `json:"room"`
	Content     string    `json:"content"`
	Timestamp   time.Time `json:"timestamp"`
	MessageType string    `json:"message_type"` // text, system, notification
}

// Room 聊天室模型
type Room struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Members     []int     `json:"members"`
	CreatedAt   time.Time `json:"created_at"`
}

// ChatServer 聊天服务器
type ChatServer struct {
	users       map[int]*User
	messages    []*Message
	rooms       map[string]*Room
	connections map[string]chan *Message // 使用通道替代WebSocket连接
	broadcast   chan *Message
	register    chan *Connection
	unregister  chan *Connection
	mutex       sync.RWMutex
	nextUserID  int
	nextMsgID   int
	eventBus    *event.EventBus
}

// Connection 连接包装
type Connection struct {
	userID int
	ch     chan *Message
}

func NewChatServer(eventBus *event.EventBus) *ChatServer {
	return &ChatServer{
		users:       make(map[int]*User),
		messages:    make([]*Message, 0),
		rooms:       make(map[string]*Room),
		connections: make(map[string]chan *Message),
		broadcast:   make(chan *Message, 100),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		nextUserID:  1,
		nextMsgID:   1,
		eventBus:    eventBus,
	}
}

// Start 启动聊天服务器
func (s *ChatServer) Start() {
	go s.run()
}

// run 运行聊天服务器的主要循环
func (s *ChatServer) run() {
	for {
		select {
		case conn := <-s.register:
			s.mutex.Lock()
			s.connections[fmt.Sprintf("%d", conn.userID)] = conn.ch
			s.updateUserOnlineStatus(conn.userID, true)
			s.mutex.Unlock()

			// 广播用户上线消息
			systemMsg := &Message{
				ID:          s.nextMsgID,
				SenderID:    0,
				Sender:      "System",
				Room:        "general",
				Content:     fmt.Sprintf("%s joined the chat", s.getUserByID(conn.userID).Username),
				Timestamp:   time.Now(),
				MessageType: "system",
			}
			s.nextMsgID++
			s.broadcast <- systemMsg

		case conn := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.connections[fmt.Sprintf("%d", conn.userID)]; ok {
				delete(s.connections, fmt.Sprintf("%d", conn.userID))
				s.updateUserOnlineStatus(conn.userID, false)
				close(conn.ch)
			}
			s.mutex.Unlock()

			// 广播用户离线消息
			systemMsg := &Message{
				ID:          s.nextMsgID,
				SenderID:    0,
				Sender:      "System",
				Room:        "general",
				Content:     fmt.Sprintf("%s left the chat", s.getUserByID(conn.userID).Username),
				Timestamp:   time.Now(),
				MessageType: "system",
			}
			s.nextMsgID++
			s.broadcast <- systemMsg

		case message := <-s.broadcast:
			s.mutex.Lock()
			// 保存消息到历史记录
			s.messages = append(s.messages, message)

			// 发送到所有连接的客户端
			for _, ch := range s.connections {
				select {
				case ch <- message:
				default:
					// 如果通道阻塞，跳过
				}
			}
			s.mutex.Unlock()
		}
	}
}

// RegisterUser 注册新用户
func (s *ChatServer) RegisterUser(username string) *User {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user := &User{
		ID:       s.nextUserID,
		Username: username,
		Online:   false,
		JoinedAt: time.Now(),
	}
	s.users[s.nextUserID] = user
	s.nextUserID++

	return user
}

// GetUserByID 根据ID获取用户
func (s *ChatServer) GetUserByID(userID int) *User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.users[userID]
}

// updateUserOnlineStatus 更新用户在线状态
func (s *ChatServer) updateUserOnlineStatus(userID int, online bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if user, exists := s.users[userID]; exists {
		user.Online = online
	}
}

// getUserByID 获取用户（私有方法，无需锁，应在已锁定的情况下调用）
func (s *ChatServer) getUserByID(userID int) *User {
	return s.users[userID]
}

// SendMessage 发送消息
func (s *ChatServer) SendMessage(senderID int, room, content string) error {
	s.mutex.RLock()
	sender := s.users[senderID]
	s.mutex.RUnlock()

	if sender == nil {
		return fmt.Errorf("user not found")
	}

	message := &Message{
		ID:          s.nextMsgID,
		SenderID:    senderID,
		Sender:      sender.Username,
		Room:        room,
		Content:     content,
		Timestamp:   time.Now(),
		MessageType: "text",
	}
	s.nextMsgID++

	s.broadcast <- message

	return nil
}

// GetRecentMessages 获取最近的消息
func (s *ChatServer) GetRecentMessages(count int) []*Message {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if count > len(s.messages) {
		count = len(s.messages)
	}

	start := len(s.messages) - count
	if start < 0 {
		start = 0
	}

	return s.messages[start:]
}

// GetOnlineUsers 获取在线用户
func (s *ChatServer) GetOnlineUsers() []*User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var onlineUsers []*User
	for _, user := range s.users {
		if user.Online {
			onlineUsers = append(onlineUsers, user)
		}
	}

	return onlineUsers
}

// CreateRoom 创建聊天室
func (s *ChatServer) CreateRoom(name, description string) *Room {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	room := &Room{
		ID:          name,
		Name:        name,
		Description: description,
		Members:     make([]int, 0),
		CreatedAt:   time.Now(),
	}

	s.rooms[name] = room
	return room
}

// JoinRoom 加入聊天室
func (s *ChatServer) JoinRoom(userID int, roomID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	room, exists := s.rooms[roomID]
	if !exists {
		return fmt.Errorf("room not found")
	}

	// 检查用户是否已在房间中
	for _, memberID := range room.Members {
		if memberID == userID {
			return nil // 用户已在房间中
		}
	}

	room.Members = append(room.Members, userID)
	return nil
}

// UserService 用户服务
type UserService struct {
	chatServer *ChatServer
}

func NewUserService(chatServer *ChatServer) *UserService {
	return &UserService{chatServer: chatServer}
}

func (s *UserService) RegisterUser(username string) *User {
	return s.chatServer.RegisterUser(username)
}

func (s *UserService) GetUser(userID int) *User {
	return s.chatServer.GetUserByID(userID)
}

func (s *UserService) GetOnlineUsers() []*User {
	return s.chatServer.GetOnlineUsers()
}

// MessageService 消息服务
type MessageService struct {
	chatServer *ChatServer
}

func NewMessageService(chatServer *ChatServer) *MessageService {
	return &MessageService{chatServer: chatServer}
}

func (s *MessageService) SendMessage(senderID int, room, content string) error {
	return s.chatServer.SendMessage(senderID, room, content)
}

func (s *MessageService) GetRecentMessages(count int) []*Message {
	return s.chatServer.GetRecentMessages(count)
}

// RoomService 房间服务
type RoomService struct {
	chatServer *ChatServer
}

func NewRoomService(chatServer *ChatServer) *RoomService {
	return &RoomService{chatServer: chatServer}
}

func (s *RoomService) CreateRoom(name, description string) *Room {
	return s.chatServer.CreateRoom(name, description)
}

func (s *RoomService) JoinRoom(userID int, roomID string) error {
	return s.chatServer.JoinRoom(userID, roomID)
}

// UserJoinedEvent 用户加入事件
type UserJoinedEvent struct {
	UserID    int
	Username  string
	Timestamp time.Time
}

func (e *UserJoinedEvent) EventName() string {
	return "user.joined"
}

// UserLeftEvent 用户离开事件
type UserLeftEvent struct {
	UserID    int
	Username  string
	Timestamp time.Time
}

func (e *UserLeftEvent) EventName() string {
	return "user.left"
}

func main() {
	// 创建容器
	container := core.New()

	// 创建事件总线
	eventBus := event.NewEventBus()

	// 创建聊天服务器
	chatServer := NewChatServer(eventBus)
	chatServer.Start()

	// 创建服务
	userSvc := NewUserService(chatServer)
	messageSvc := NewMessageService(chatServer)
	roomSvc := NewRoomService(chatServer)

	// 注册服务到容器
	container.Register("chatServer", core.Bean(chatServer))
	container.Register("userService", core.Bean(userSvc))
	container.Register("messageService", core.Bean(messageSvc))
	container.Register("roomService", core.Bean(roomSvc))
	container.Register("eventBus", core.Bean(eventBus))

	// 创建Gin实例并绑定容器
	g := ginadapter.New(ginadapter.WithContainer(container))

	// 设置路由
	g.Engine().LoadHTMLGlob("templates/*")

	g.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "实时聊天系统",
		})
	})

	// REST API路由
	g.POST("/api/users", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
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
		user := userSvc.RegisterUser(req.Username)

		c.JSON(http.StatusOK, user)
	})

	g.GET("/api/users/:id", func(c *gin.Context) {
		// 这里应该是从路径参数获取用户ID的实际实现
		// 为了简化示例，我们返回默认用户
		container := g.Container().(core.Container)
		userSvcRaw, err := container.Get("userService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userSvc := userSvcRaw.(*UserService)
		user := userSvc.GetUser(1) // 简化实现

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	g.GET("/api/users/online", func(c *gin.Context) {
		container := g.Container().(core.Container)
		userSvcRaw, err := container.Get("userService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userSvc := userSvcRaw.(*UserService)
		onlineUsers := userSvc.GetOnlineUsers()

		c.JSON(http.StatusOK, onlineUsers)
	})

	g.POST("/api/messages", func(c *gin.Context) {
		var req struct {
			SenderID int    `json:"sender_id"`
			Room     string `json:"room"`
			Content  string `json:"content"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		container := g.Container().(core.Container)
		messageSvcRaw, err := container.Get("messageService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		messageSvc := messageSvcRaw.(*MessageService)
		err = messageSvc.SendMessage(req.SenderID, req.Room, req.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Message sent"})
	})

	g.GET("/api/messages/recent", func(c *gin.Context) {
		count := 50 // 默认获取50条消息
		countStr := c.Query("count")
		if countStr != "" {
			if c, err := fmt.Sscanf(countStr, "%d", &count); c > 0 && err == nil {
				if count > 100 {
					count = 100 // 限制最大数量
				}
			}
		}

		container := g.Container().(core.Container)
		messageSvcRaw, err := container.Get("messageService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		messageSvc := messageSvcRaw.(*MessageService)
		messages := messageSvc.GetRecentMessages(count)

		c.JSON(http.StatusOK, messages)
	})

	g.POST("/api/rooms", func(c *gin.Context) {
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		container := g.Container().(core.Container)
		roomSvcRaw, err := container.Get("roomService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		roomSvc := roomSvcRaw.(*RoomService)
		room := roomSvc.CreateRoom(req.Name, req.Description)

		c.JSON(http.StatusOK, room)
	})

	g.POST("/api/rooms/:id/join", func(c *gin.Context) {
		var req struct {
			UserID int `json:"user_id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		roomID := c.Param("id")
		container := g.Container().(core.Container)
		roomSvcRaw, err := container.Get("roomService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		roomSvc := roomSvcRaw.(*RoomService)
		err = roomSvc.JoinRoom(req.UserID, roomID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Joined room successfully"})
	})

	fmt.Println("实时聊天系统启动在 :8085")

	// 启动服务器
	if err := g.Start(); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}
