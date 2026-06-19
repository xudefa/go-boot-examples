// Package main 演示 go-boot 的安全模块和 JWT 模块集成
//
// 本示例展示如何:
//   - 同时启用安全模块和 JWT 模块
//   - 配置用户信息（内存存储）
//   - 登录获取 JWT Token
//   - 使用 Token 访问受保护的 API
//   - 访问排除路径（不需要认证）
//
// 使用方式:
//
//	cd examples/jwt-security
//	go run .
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/xudefa/go-boot/boot"
	"github.com/xudefa/go-boot/core"
	"github.com/xudefa/go-boot-jwt"
	"github.com/xudefa/go-boot/security"
)

var app *boot.Boot

func main() {
	fmt.Println("=== JWT Security Integration Example ===")
	fmt.Println()

	var err error
	app, err = boot.NewApplication(
		boot.WithAppName("jwt-security-example"),
		boot.WithVersion("1.0.0"),
		boot.WithProfiles("dev"),
		boot.WithConfigLocation("./application.json"),
	)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// 注册自定义用户服务（添加测试用户）
	uds := security.NewInMemoryUserDetailsService()
	uds.CreateUser("admin", "admin123", []string{"ROLE_ADMIN", "ROLE_USER"})
	uds.CreateUser("user", "user123", []string{"ROLE_USER"})
	app.Container().Register("customUserDetailsService",
		core.Bean(uds),
		core.Singleton(),
	)

	// 注册 NoOp 密码编码器（用于开发测试，明文存储密码）
	app.Container().Register("passwordEncoder",
		core.Bean(security.NewNoOpPasswordEncoder()),
		core.Singleton(),
	)

	// 启动应用生命周期（初始化自动配置）
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	fmt.Println("Application started successfully!")
	fmt.Println("Server running on http://localhost:8080")
	fmt.Println()
	fmt.Println("Endpoints:")
	fmt.Println("  POST   /login          - 登录获取 Token（不需要认证）")
	fmt.Println("  GET    /health         - 健康检查（不需要认证）")
	fmt.Println("  GET    /api/public     - 公开 API（不需要认证）")
	fmt.Println("  GET    /api/protected  - 受保护 API（需要 JWT Token）")
	fmt.Println("  GET    /api/admin      - 管理员 API（需要 JWT Token 和 ADMIN 角色）")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  1. 登录获取 Token:")
	fmt.Println("     curl -X POST http://localhost:8080/login -H 'Content-Type: application/json' -d '{\"username\":\"admin\",\"password\":\"admin123\"}'")
	fmt.Println()
	fmt.Println("  2. 使用 Token 访问受保护 API:")
	fmt.Println("     curl http://localhost:8080/api/protected -H 'Authorization: Bearer <token>'")
	fmt.Println()

	// 获取安全过滤器链
	securityFilterChain, err := app.Container().Get("securityFilterChain")
	if err != nil {
		log.Printf("Warning: Security filter chain not found: %v", err)
		// 如果没有安全过滤器链，使用默认的 HTTP 服务器
		log.Fatal(http.ListenAndServe(":8080", nil))
		return
	}

	// 使用安全过滤器链包装 HTTP 服务器
	if filterChain, ok := securityFilterChain.(security.SecurityFilterChain); ok {
		securityHandler := security.NewSecurityFilterChainHandler(filterChain, http.DefaultServeMux)
		fmt.Println("Security filter chain integrated successfully!")
		// 启动带有安全过滤器的 HTTP 服务器
		log.Fatal(http.ListenAndServe(":8080", securityHandler))
		return
	}

	// 如果类型断言失败，使用默认的 HTTP 服务器
	fmt.Printf("Warning: Security filter chain type assertion failed")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// LoginHandler 登录处理
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if loginRequest.Username == "" || loginRequest.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	authManagerBean, err := app.Container().Get("authenticationManager")
	if err != nil {
		http.Error(w, "AuthenticationManager not found", http.StatusInternalServerError)
		return
	}
	authManager, ok := authManagerBean.(security.AuthenticationManager)
	if !ok {
		http.Error(w, "AuthenticationManager type assertion failed", http.StatusInternalServerError)
		return
	}

	authToken := security.NewUsernamePasswordAuthenticationToken(loginRequest.Username, loginRequest.Password)
	authenticated, err := authManager.Authenticate(r.Context(), authToken)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	bean, err := app.Container().Get("jwtUtil")
	if err != nil {
		http.Error(w, "JwtUtil not found", http.StatusInternalServerError)
		return
	}
	jwtUtil, ok := bean.(*jwt.JwtUtil)
	if !ok {
		http.Error(w, "JwtUtil type assertion failed", http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := jwtUtil.GenerateToken(authenticated.Name(), "web-app")
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// PublicHandler 公开接口（不需要认证）
func PublicHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "This is a public API, no authentication required",
	})
}

// ProtectedHandler 受保护接口（需要 JWT Token）
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	auth := security.GetAuthenticationFromContext(r.Context())
	if auth == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Not authenticated",
		})
		return
	}

	response := map[string]interface{}{
		"message":       "This is a protected API",
		"authenticated": auth.Authenticated(),
		"username":      auth.Name(),
		"authorities":   auth.Authorities(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// AdminHandler 管理员接口（需要 JWT Token 和 ADMIN 角色）
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	auth := security.GetAuthenticationFromContext(r.Context())
	if auth == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Not authenticated",
		})
		return
	}

	response := map[string]interface{}{
		"message":       "This is an admin API",
		"authenticated": auth.Authenticated(),
		"username":      auth.Name(),
		"authorities":   auth.Authorities(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// HealthHandler 健康检查接口（不需要认证）
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "UP",
	})
}

func init() {
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/api/public", PublicHandler)
	http.HandleFunc("/api/protected", ProtectedHandler)
	http.HandleFunc("/api/admin", AdminHandler)
}
