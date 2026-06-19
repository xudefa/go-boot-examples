package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xudefa/go-boot/core"
	"github.com/xudefa/go-boot/event"
	ginadapter "github.com/xudefa/go-boot-gin/server"
	"github.com/xudefa/go-boot/schedule"
)

// User 用户模型
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// Task 任务模型
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AssigneeID  int       `json:"assignee_id"`
	CreatorID   int       `json:"creator_id"`
	Status      string    `json:"status"`   // pending, in_progress, completed, cancelled
	Priority    string    `json:"priority"` // low, medium, high, urgent
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskReport 任务报告模型
type TaskReport struct {
	ID          int       `json:"id"`
	TaskID      int       `json:"task_id"`
	GeneratedAt time.Time `json:"generated_at"`
	Content     string    `json:"content"`
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

// TaskService 任务服务
type TaskService struct {
	tasks  map[int]*Task
	nextID int
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}
}

func (s *TaskService) CreateTask(title, description string, assigneeID, creatorID int, priority string, dueDate time.Time) *Task {
	task := &Task{
		ID:          s.nextID,
		Title:       title,
		Description: description,
		AssigneeID:  assigneeID,
		CreatorID:   creatorID,
		Status:      "pending",
		Priority:    priority,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.tasks[s.nextID] = task
	s.nextID++
	return task
}

func (s *TaskService) GetTask(id int) *Task {
	return s.tasks[id]
}

func (s *TaskService) UpdateTaskStatus(id int, status string) *Task {
	task, exists := s.tasks[id]
	if !exists {
		return nil
	}

	task.Status = status
	task.UpdatedAt = time.Now()
	return task
}

func (s *TaskService) UpdateTask(id int, title, description, status, priority string, dueDate time.Time) *Task {
	task, exists := s.tasks[id]
	if !exists {
		return nil
	}

	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	if status != "" {
		task.Status = status
	}
	if priority != "" {
		task.Priority = priority
	}
	if !dueDate.IsZero() {
		task.DueDate = dueDate
	}

	task.UpdatedAt = time.Now()
	return task
}

func (s *TaskService) GetTasksByAssignee(assigneeID int) []*Task {
	var tasks []*Task
	for _, task := range s.tasks {
		if task.AssigneeID == assigneeID {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (s *TaskService) GetTasksByStatus(status string) []*Task {
	var tasks []*Task
	for _, task := range s.tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (s *TaskService) GetOverdueTasks() []*Task {
	var overdueTasks []*Task
	now := time.Now()
	for _, task := range s.tasks {
		if task.Status != "completed" && task.Status != "cancelled" && task.DueDate.Before(now) {
			overdueTasks = append(overdueTasks, task)
		}
	}
	return overdueTasks
}

// ReminderService 提醒服务
type ReminderService struct {
	taskSvc  *TaskService
	eventBus *event.EventBus
}

func NewReminderService(taskSvc *TaskService, eventBus *event.EventBus) *ReminderService {
	return &ReminderService{
		taskSvc:  taskSvc,
		eventBus: eventBus,
	}
}

// CheckDueTasks 检查即将到期的任务
func (s *ReminderService) CheckDueTasks() {
	now := time.Now()
	nextHour := now.Add(time.Hour)

	for _, task := range s.taskSvc.tasks {
		if task.Status != "completed" && task.Status != "cancelled" &&
			task.DueDate.After(now) && task.DueDate.Before(nextHour) {
			// 发送即将到期提醒事件
			reminderEvent := &TaskReminderEvent{
				TaskID:    task.ID,
				TaskTitle: task.Title,
				DueDate:   task.DueDate,
				EventType: "due_soon",
				EventTime: time.Now(),
			}
			s.eventBus.Publish(reminderEvent)
		}
	}
}

// CheckOverdueTasks 检查逾期任务
func (s *ReminderService) CheckOverdueTasks() {
	for _, task := range s.taskSvc.GetOverdueTasks() {
		// 发送逾期提醒事件
		reminderEvent := &TaskReminderEvent{
			TaskID:    task.ID,
			TaskTitle: task.Title,
			DueDate:   task.DueDate,
			EventType: "overdue",
			EventTime: time.Now(),
		}
		s.eventBus.Publish(reminderEvent)
	}
}

// TaskReportService 报告服务
type TaskReportService struct {
	reports map[int]*TaskReport
	nextID  int
	taskSvc *TaskService
}

func NewTaskReportService(taskSvc *TaskService) *TaskReportService {
	return &TaskReportService{
		reports: make(map[int]*TaskReport),
		nextID:  1,
		taskSvc: taskSvc,
	}
}

func (s *TaskReportService) GenerateDailyReport() *TaskReport {
	now := time.Now()

	// 统计今天的任务完成情况
	var completedTasks, pendingTasks, overdueTasks int

	for _, task := range s.taskSvc.tasks {
		if task.CreatedAt.Day() == now.Day() && task.CreatedAt.Month() == now.Month() {
			if task.Status == "completed" {
				completedTasks++
			} else if task.Status == "pending" {
				pendingTasks++
			}
		}

		if task.Status != "completed" && task.Status != "cancelled" && task.DueDate.Before(now) {
			overdueTasks++
		}
	}

	content := fmt.Sprintf("任务日报 - %s\n\n今日新增任务: %d\n已完成任务: %d\n待处理任务: %d\n逾期任务: %d",
		now.Format("2006-01-02"), 0, completedTasks, pendingTasks, overdueTasks)

	report := &TaskReport{
		ID:          s.nextID,
		TaskID:      0, // 日报不关联特定任务
		GeneratedAt: now,
		Content:     content,
	}

	s.reports[s.nextID] = report
	s.nextID++
	return report
}

func (s *TaskReportService) GenerateUserReport(userID int) *TaskReport {
	userTasks := s.taskSvc.GetTasksByAssignee(userID)

	var completed, pending, inProgress int
	for _, task := range userTasks {
		switch task.Status {
		case "completed":
			completed++
		case "pending":
			pending++
		case "in_progress":
			inProgress++
		}
	}

	now := time.Now()
	content := fmt.Sprintf("用户任务报告 - %s\n\n分配给您的任务总数: %d\n已完成: %d\n进行中: %d\n待处理: %d",
		now.Format("2006-01-02"), len(userTasks), completed, inProgress, pending)

	report := &TaskReport{
		ID:          s.nextID,
		TaskID:      0, // 用户报告不关联特定任务
		GeneratedAt: now,
		Content:     content,
	}

	s.reports[s.nextID] = report
	s.nextID++
	return report
}

// TaskCreatedEvent 任务创建事件
type TaskCreatedEvent struct {
	TaskID    int
	TaskTitle string
	Assignee  string
	EventTime time.Time
}

func (e *TaskCreatedEvent) Type() string {
	return "task.created"
}

func (e *TaskCreatedEvent) Timestamp() time.Time {
	return e.EventTime
}

// TaskStatusChangedEvent 任务状态变更事件
type TaskStatusChangedEvent struct {
	TaskID    int
	TaskTitle string
	OldStatus string
	NewStatus string
	EventTime time.Time
	ChangedBy string
}

func (e *TaskStatusChangedEvent) Type() string {
	return "task.status.changed"
}

func (e *TaskStatusChangedEvent) Timestamp() time.Time {
	return e.EventTime
}

// TaskReminderEvent 任务提醒事件
type TaskReminderEvent struct {
	TaskID    int
	TaskTitle string
	DueDate   time.Time
	EventType string // due_soon, overdue
	EventTime time.Time
}

func (e *TaskReminderEvent) Type() string {
	return "task.reminder"
}

func (e *TaskReminderEvent) GetEventType() string {
	return e.EventType
}

func (e *TaskReminderEvent) Timestamp() time.Time {
	return e.EventTime
}

func main() {
	// 创建容器
	container := core.New()

	// 创建事件总线
	eventBus := event.NewEventBus()

	// 创建服务
	userSvc := NewUserService()
	taskSvc := NewTaskService()
	reminderSvc := NewReminderService(taskSvc, eventBus)
	reportSvc := NewTaskReportService(taskSvc)

	// 注册服务到容器
	container.Register("userService", core.Bean(userSvc))
	container.Register("taskService", core.Bean(taskSvc))
	container.Register("reminderService", core.Bean(reminderSvc))
	container.Register("reportService", core.Bean(reportSvc))
	container.Register("eventBus", core.Bean(eventBus))

	// 订阅事件
	eventBus.Subscribe("task.reminder", func(e event.ApplicationEvent) {
		if reminderEvent, ok := e.(*TaskReminderEvent); ok {
			fmt.Printf("[REMINDER] Task '%s' (ID: %d) is %s at %s\n",
				reminderEvent.TaskTitle, reminderEvent.TaskID, reminderEvent.GetEventType(), reminderEvent.DueDate.Format("2006-01-02 15:04:05"))
		}
	})

	// 创建Gin实例并绑定容器
	g := ginadapter.New(ginadapter.WithContainer(container))

	// 设置路由
	setupRoutes(g)

	// 启动定时任务
	go startScheduler(reminderSvc, reportSvc)

	fmt.Println("任务管理系统启动在 :8082")

	// 启动服务器
	if err := g.Start(); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}

func setupRoutes(g *ginadapter.GinServer) {
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
			req.Role = "user"
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

	// 任务相关路由
	g.GET("/api/tasks", func(c *gin.Context) {
		status := c.Query("status")

		container := g.Container().(core.Container)
		taskSvcRaw, err := container.Get("taskService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		taskSvc := taskSvcRaw.(*TaskService)
		var tasks []*Task

		if status != "" {
			tasks = taskSvc.GetTasksByStatus(status)
		} else {
			// 返回所有任务
			tasks = make([]*Task, 0, len(taskSvc.tasks))
			for _, task := range taskSvc.tasks {
				tasks = append(tasks, task)
			}
		}

		c.JSON(http.StatusOK, tasks)
	})

	g.GET("/api/tasks/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		container := g.Container().(core.Container)
		taskSvcRaw, err := container.Get("taskService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		taskSvc := taskSvcRaw.(*TaskService)
		task := taskSvc.GetTask(id)
		if task == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		c.JSON(http.StatusOK, task)
	})

	g.POST("/api/tasks", func(c *gin.Context) {
		var req struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			AssigneeID  int       `json:"assignee_id"`
			CreatorID   int       `json:"creator_id"`
			Priority    string    `json:"priority"`
			DueDate     time.Time `json:"due_date"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Priority == "" {
			req.Priority = "medium"
		}

		container := g.Container().(core.Container)
		taskSvcRaw, err := container.Get("taskService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		taskSvc := taskSvcRaw.(*TaskService)
		task := taskSvc.CreateTask(req.Title, req.Description, req.AssigneeID, req.CreatorID, req.Priority, req.DueDate)

		// 发布任务创建事件
		eventBusRaw, err := g.Container().(core.Container).Get("eventBus")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		eventBus := eventBusRaw.(*event.EventBus)

		userSvcRaw, err := g.Container().(core.Container).Get("userService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userSvc := userSvcRaw.(*UserService)

		assignee := userSvc.GetUser(req.AssigneeID)
		assigneeName := "Unknown"
		if assignee != nil {
			assigneeName = assignee.Username
		}

		taskCreatedEvent := &TaskCreatedEvent{
			TaskID:    task.ID,
			TaskTitle: task.Title,
			Assignee:  assigneeName,
			EventTime: time.Now(),
		}
		eventBus.Publish(taskCreatedEvent)

		c.JSON(http.StatusOK, task)
	})

	g.PUT("/api/tasks/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		var req struct {
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Status      string    `json:"status"`
			Priority    string    `json:"priority"`
			DueDate     time.Time `json:"due_date"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		container := g.Container().(core.Container)
		taskSvcRaw, err := container.Get("taskService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		taskSvc := taskSvcRaw.(*TaskService)

		// 保存旧状态用于事件
		oldTask := taskSvc.GetTask(id)
		if oldTask == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		oldStatus := oldTask.Status

		task := taskSvc.UpdateTask(id, req.Title, req.Description, req.Status, req.Priority, req.DueDate)
		if task == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		// 如果状态改变了，发布状态变更事件
		if oldStatus != task.Status {
			container := g.Container().(core.Container)
			eventBusRaw, err := container.Get("eventBus")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			eventBus := eventBusRaw.(*event.EventBus)

			userSvcRaw, err := container.Get("userService")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			userSvc := userSvcRaw.(*UserService)

			changedBy := "System"
			if task.AssigneeID != 0 {
				if user := userSvc.GetUser(task.AssigneeID); user != nil {
					changedBy = user.Username
				}
			}

			statusChangedEvent := &TaskStatusChangedEvent{
				TaskID:    task.ID,
				TaskTitle: task.Title,
				OldStatus: oldStatus,
				NewStatus: task.Status,
				EventTime: time.Now(),
				ChangedBy: changedBy,
			}
			eventBus.Publish(statusChangedEvent)
		}

		c.JSON(http.StatusOK, task)
	})

	g.PUT("/api/tasks/:id/status", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
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
		taskSvcRaw, err := container.Get("taskService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		taskSvc := taskSvcRaw.(*TaskService)

		// 保存旧状态用于事件
		oldTask := taskSvc.GetTask(id)
		if oldTask == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		oldStatus := oldTask.Status

		task := taskSvc.UpdateTaskStatus(id, req.Status)
		if task == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		// 发布状态变更事件
		eventBusRaw, err := container.Get("eventBus")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		eventBus := eventBusRaw.(*event.EventBus)

		userSvcRaw, err := container.Get("userService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userSvc := userSvcRaw.(*UserService)

		changedBy := "System"
		if task.AssigneeID != 0 {
			if user := userSvc.GetUser(task.AssigneeID); user != nil {
				changedBy = user.Username
			}
		}

		statusChangedEvent := &TaskStatusChangedEvent{
			TaskID:    task.ID,
			TaskTitle: task.Title,
			OldStatus: oldStatus,
			NewStatus: task.Status,
			EventTime: time.Now(),
			ChangedBy: changedBy,
		}
		eventBus.Publish(statusChangedEvent)

		c.JSON(http.StatusOK, task)
	})

	g.GET("/api/users/:id/tasks", func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		container := g.Container().(core.Container)
		taskSvcRaw, err := container.Get("taskService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		taskSvc := taskSvcRaw.(*TaskService)
		tasks := taskSvc.GetTasksByAssignee(userId)

		c.JSON(http.StatusOK, tasks)
	})

	g.GET("/api/tasks/overdue", func(c *gin.Context) {
		container := g.Container().(core.Container)
		taskSvcRaw, err := container.Get("taskService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		taskSvc := taskSvcRaw.(*TaskService)
		overdueTasks := taskSvc.GetOverdueTasks()

		c.JSON(http.StatusOK, overdueTasks)
	})

	// 报告相关路由
	g.GET("/api/reports/daily", func(c *gin.Context) {
		container := g.Container().(core.Container)
		reportSvcRaw, err := container.Get("reportService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		reportSvc := reportSvcRaw.(*TaskReportService)
		report := reportSvc.GenerateDailyReport()

		c.JSON(http.StatusOK, report)
	})

	g.GET("/api/reports/user/:id", func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		container := g.Container().(core.Container)
		reportSvcRaw, err := container.Get("reportService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		reportSvc := reportSvcRaw.(*TaskReportService)
		report := reportSvc.GenerateUserReport(userId)

		c.JSON(http.StatusOK, report)
	})
}

func startScheduler(reminderSvc *ReminderService, reportSvc *TaskReportService) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := schedule.NewScheduler()

	// 每5分钟检查一次即将到期的任务
	task1 := schedule.NewTask("*/5 * * * * ?", "check_due_tasks", func(ctx context.Context) error {
		fmt.Println("Checking due tasks...")
		reminderSvc.CheckDueTasks()
		return nil
	})
	scheduler.Register(task1)

	// 每小时检查一次逾期任务
	task2 := schedule.NewTask("0 * * * * ?", "check_overdue_tasks", func(ctx context.Context) error {
		fmt.Println("Checking overdue tasks...")
		reminderSvc.CheckOverdueTasks()
		return nil
	})
	scheduler.Register(task2)

	// 每天上午9点生成日报
	task3 := schedule.NewTask("0 9 * * * ?", "generate_daily_report", func(ctx context.Context) error {
		fmt.Println("Generating daily report...")
		report := reportSvc.GenerateDailyReport()
		fmt.Printf("Daily report generated: %s\n", report.Content)
		return nil
	})
	scheduler.Register(task3)

	scheduler.Start(ctx)
}
