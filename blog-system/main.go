package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

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
}

// Article 文章模型
type Article struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Author  User      `json:"author"`
	Created time.Time `json:"created"`
}

// Comment 评论模型
type Comment struct {
	ID        int       `json:"id"`
	ArticleID int       `json:"article_id"`
	Content   string    `json:"content"`
	Author    User      `json:"author"`
	Created   time.Time `json:"created"`
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

func (s *UserService) CreateUser(username, email string) *User {
	user := &User{
		ID:       s.nextID,
		Username: username,
		Email:    email,
	}
	s.users[s.nextID] = user
	s.nextID++
	return user
}

func (s *UserService) GetUser(id int) *User {
	return s.users[id]
}

// ArticleService 文章服务
type ArticleService struct {
	articles map[int]*Article
	nextID   int
	userSvc  *UserService
}

func NewArticleService(userSvc *UserService) *ArticleService {
	return &ArticleService{
		articles: make(map[int]*Article),
		nextID:   1,
		userSvc:  userSvc,
	}
}

func (s *ArticleService) CreateArticle(title, content string, authorID int) *Article {
	author := s.userSvc.GetUser(authorID)
	if author == nil {
		return nil // 用户不存在
	}

	article := &Article{
		ID:      s.nextID,
		Title:   title,
		Content: content,
		Author:  *author,
		Created: time.Now(),
	}
	s.articles[s.nextID] = article
	s.nextID++
	return article
}

func (s *ArticleService) GetArticle(id int) *Article {
	return s.articles[id]
}

func (s *ArticleService) GetAllArticles() []*Article {
	articles := make([]*Article, 0, len(s.articles))
	for _, article := range s.articles {
		articles = append(articles, article)
	}
	return articles
}

func (s *ArticleService) UpdateArticle(id int, title, content string) *Article {
	article, exists := s.articles[id]
	if !exists {
		return nil
	}

	article.Title = title
	article.Content = content
	return article
}

func (s *ArticleService) DeleteArticle(id int) bool {
	_, exists := s.articles[id]
	if exists {
		delete(s.articles, id)
		return true
	}
	return false
}

// CommentService 评论服务
type CommentService struct {
	comments   map[int]*Comment
	nextID     int
	userSvc    *UserService
	articleSvc *ArticleService
}

func NewCommentService(userSvc *UserService, articleSvc *ArticleService) *CommentService {
	return &CommentService{
		comments:   make(map[int]*Comment),
		nextID:     1,
		userSvc:    userSvc,
		articleSvc: articleSvc,
	}
}

func (s *CommentService) CreateComment(articleID int, content string, authorID int) *Comment {
	article := s.articleSvc.GetArticle(articleID)
	if article == nil {
		return nil // 文章不存在
	}

	author := s.userSvc.GetUser(authorID)
	if author == nil {
		return nil // 用户不存在
	}

	comment := &Comment{
		ID:        s.nextID,
		ArticleID: articleID,
		Content:   content,
		Author:    *author,
		Created:   time.Now(),
	}
	s.comments[s.nextID] = comment
	s.nextID++
	return comment
}

func (s *CommentService) GetCommentsByArticle(articleID int) []*Comment {
	comments := make([]*Comment, 0)
	for _, comment := range s.comments {
		if comment.ArticleID == articleID {
			comments = append(comments, comment)
		}
	}
	return comments
}

// BlogEventService 事件服务
type BlogEventService struct {
	eventBus *event.EventBus
}

func NewBlogEventService(eventBus *event.EventBus) *BlogEventService {
	return &BlogEventService{eventBus: eventBus}
}

// ArticleCreatedEvent 文章创建事件
type ArticleCreatedEvent struct {
	EventType string
	EventTime time.Time
	ArticleID int
	Title     string
}

func (e *ArticleCreatedEvent) Type() string {
	return e.EventType
}

func (e *ArticleCreatedEvent) Timestamp() time.Time {
	return e.EventTime
}

// PublishArticleCreatedEvent 发布文章创建事件
func (s *BlogEventService) PublishArticleCreatedEvent(articleID int, title string) {
	event := &ArticleCreatedEvent{
		EventType: "article.created",
		EventTime: time.Now(),
		ArticleID: articleID,
		Title:     title,
	}
	s.eventBus.Publish(event)
}

func main() {
	// 创建容器
	container := core.New()

	// 创建事件总线
	eventBus := event.NewEventBus()

	// 创建服务
	userSvc := NewUserService()
	articleSvc := NewArticleService(userSvc)
	commentSvc := NewCommentService(userSvc, articleSvc)
	blogEventSvc := NewBlogEventService(eventBus)

	// 注册服务到容器
	container.Register("userService", core.Bean(userSvc))
	container.Register("articleService", core.Bean(articleSvc))
	container.Register("commentService", core.Bean(commentSvc))
	container.Register("blogEventService", core.Bean(blogEventSvc))
	container.Register("eventBus", core.Bean(eventBus))

	// 创建Gin实例并绑定容器
	g := ginadapter.New(ginadapter.WithContainer(container))

	// 注册路由
	setupRoutes(g)

	fmt.Println("博客系统启动在 :8080")

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
		user := userSvc.CreateUser(req.Username, req.Email)

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

	// 文章相关路由
	g.GET("/api/articles", func(c *gin.Context) {
		container := g.Container().(core.Container)
		articleSvcRaw, err := container.Get("articleService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		articleSvc := articleSvcRaw.(*ArticleService)
		articles := articleSvc.GetAllArticles()

		c.JSON(http.StatusOK, articles)
	})

	g.GET("/api/articles/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
			return
		}

		container := g.Container().(core.Container)
		articleSvcRaw, err := container.Get("articleService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		articleSvc := articleSvcRaw.(*ArticleService)
		article := articleSvc.GetArticle(id)
		if article == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}

		c.JSON(http.StatusOK, article)
	})

	g.POST("/api/articles", func(c *gin.Context) {
		var req struct {
			Title    string `json:"title"`
			Content  string `json:"content"`
			AuthorID int    `json:"author_id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		container := g.Container().(core.Container)
		articleSvcRaw, err := container.Get("articleService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		articleSvc := articleSvcRaw.(*ArticleService)
		article := articleSvc.CreateArticle(req.Title, req.Content, req.AuthorID)
		if article == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create article"})
			return
		}

		// 发布文章创建事件
		eventSvcRaw, err := container.Get("blogEventService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		eventSvc := eventSvcRaw.(*BlogEventService)
		eventSvc.PublishArticleCreatedEvent(article.ID, article.Title)

		c.JSON(http.StatusOK, article)
	})

	g.PUT("/api/articles/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
			return
		}

		var req struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		container := g.Container().(core.Container)
		articleSvcRaw, err := container.Get("articleService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		articleSvc := articleSvcRaw.(*ArticleService)
		updated := articleSvc.UpdateArticle(id, req.Title, req.Content)
		if updated == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}

		c.JSON(http.StatusOK, updated)
	})

	g.DELETE("/api/articles/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
			return
		}

		container := g.Container().(core.Container)
		articleSvcRaw, err := container.Get("articleService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		articleSvc := articleSvcRaw.(*ArticleService)
		success := articleSvc.DeleteArticle(id)
		if !success {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
	})

	// 评论相关路由
	g.GET("/api/articles/:id/comments", func(c *gin.Context) {
		articleID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
			return
		}

		container := g.Container().(core.Container)
		commentSvcRaw, err := container.Get("commentService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		commentSvc := commentSvcRaw.(*CommentService)
		comments := commentSvc.GetCommentsByArticle(articleID)

		c.JSON(http.StatusOK, comments)
	})

	g.POST("/api/comments", func(c *gin.Context) {
		var req struct {
			ArticleID int    `json:"article_id"`
			Content   string `json:"content"`
			AuthorID  int    `json:"author_id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		container := g.Container().(core.Container)
		commentSvcRaw, err := container.Get("commentService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		commentSvc := commentSvcRaw.(*CommentService)
		comment := commentSvc.CreateComment(req.ArticleID, req.Content, req.AuthorID)
		if comment == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create comment"})
			return
		}

		c.JSON(http.StatusOK, comment)
	})

	// 静态页面路由
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "templates/*")
	g.Engine().LoadHTMLGlob(templatePath)

	g.GET("/", func(c *gin.Context) {
		container := g.Container().(core.Container)
		articleSvcRaw, err := container.Get("articleService")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		articleSvc := articleSvcRaw.(*ArticleService)
		articles := articleSvc.GetAllArticles()

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":    "博客首页",
			"articles": articles,
		})
	})
}
