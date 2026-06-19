// Package main 演示 go-boot 的 JWT 工具
//
// 本示例展示如何:
//   - 创建 JWT 工具实例
//   - 生成 Token 和 RefreshToken
//   - 解析和验证 Token
//   - 刷新 Token
//   - 获取 Token 中的声明信息
//   - 获取 Token 剩余有效期
//
// 使用方式:
//
//	cd examples/jwt-example
//	go run .
package main

import (
	"fmt"
	"time"

	"github.com/xudefa/go-boot-jwt"
)

func main() {
	fmt.Println("=== JWT Example ===")
	fmt.Println()

	// 创建 JWT 工具实例，使用自定义配置
	jwtUtil := jwt.NewJWTUtil(
		jwt.WithSecretKey("my-secret-key-123456"),
		jwt.WithIssuer("go-boot-example"),
		jwt.WithExpiresDuration(30*time.Minute),
		jwt.WithRefreshExpiresDuration(2*time.Hour),
	)
	fmt.Println("JWT utility created with custom config")
	fmt.Println("  Secret Key: my-secret-key-123456")
	fmt.Println("  Issuer: go-boot-example")
	fmt.Println("  Expires Duration: 30 minutes")
	fmt.Println("  Refresh Expires Duration: 2 hours")
	fmt.Println()

	// 生成 Token 和 RefreshToken
	fmt.Println("=== Generate Token ===")
	username := "alice"
	audience := []string{"web-app", "mobile-app"}
	token, refreshToken, err := jwtUtil.GenerateToken(username, audience...)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return
	}
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Audience: %v\n", audience)
	fmt.Printf("Access Token: %s\n", token)
	fmt.Printf("Refresh Token: %s\n", refreshToken)
	fmt.Println()

	// 解析 Token
	fmt.Println("=== Parse Token ===")
	claims, err := jwtUtil.ParseToken(token)
	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err)
		return
	}
	fmt.Printf("Token ID: %s\n", claims.ID)
	fmt.Printf("Issuer: %s\n", claims.Issuer)
	fmt.Printf("Subject (Username): %s\n", claims.Subject)
	fmt.Printf("Audience: %v\n", claims.Audience)
	fmt.Printf("Issued At: %s\n", claims.IssuedAt.Time.Format(time.RFC3339))
	fmt.Printf("Expires At: %s\n", claims.ExpiresAt.Time.Format(time.RFC3339))
	fmt.Printf("Not Before: %s\n", claims.NotBefore.Time.Format(time.RFC3339))
	fmt.Println()

	// 验证 Token
	fmt.Println("=== Validate Token ===")
	valid, err := jwtUtil.ValidateToken(token)
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
	} else {
		fmt.Printf("Token is valid: %v\n", valid)
	}
	fmt.Println()

	// 获取 Token 中的用户信息
	fmt.Println("=== Get Token Info ===")
	subject, err := jwtUtil.GetSubject(token)
	if err != nil {
		fmt.Printf("Error getting subject: %v\n", err)
	} else {
		fmt.Printf("Subject (Username): %s\n", subject)
	}

	userID, err := jwtUtil.GetUserId(token)
	if err != nil {
		fmt.Printf("Error getting user ID: %v\n", err)
	} else {
		fmt.Printf("User ID: %s\n", userID)
	}

	claimsInfo, err := jwtUtil.GetClaims(token)
	if err != nil {
		fmt.Printf("Error getting claims: %v\n", err)
	} else {
		fmt.Printf("Full Claims: %+v\n", claimsInfo)
	}
	fmt.Println()

	// 获取 Token 剩余有效期
	fmt.Println("=== Get Remaining Time ===")
	remainingTime, err := jwtUtil.GetRemainingTime(token)
	if err != nil {
		fmt.Printf("Error getting remaining time: %v\n", err)
	} else {
		fmt.Printf("Remaining time: %v\n", remainingTime)
	}
	fmt.Println()

	// 刷新 Token
	fmt.Println("=== Refresh Token ===")
	newToken, newRefreshToken, err := jwtUtil.RefreshToken(token)
	if err != nil {
		fmt.Printf("Error refreshing token: %v\n", err)
		return
	}
	fmt.Printf("New Access Token: %s\n", newToken)
	fmt.Printf("New Refresh Token: %s\n", newRefreshToken)

	// 验证新 Token
	newClaims, err := jwtUtil.ParseToken(newToken)
	if err != nil {
		fmt.Printf("Error parsing new token: %v\n", err)
		return
	}
	fmt.Printf("New Token ID: %s\n", newClaims.ID)
	fmt.Printf("New Token Expires At: %s\n", newClaims.ExpiresAt.Time.Format(time.RFC3339))
	fmt.Println()

	// 测试无效 Token
	fmt.Println("=== Test Invalid Token ===")
	invalidToken := "invalid.token.string"
	valid, err = jwtUtil.ValidateToken(invalidToken)
	if err != nil {
		fmt.Printf("Expected error for invalid token: %v\n", err)
	} else {
		fmt.Printf("Unexpected result: valid=%v\n", valid)
	}
	fmt.Println()

	// 测试过期 Token（使用一个已过期的 Token）
	fmt.Println("=== Test Expired Token ===")
	expiredJwtUtil := jwt.NewJWTUtil(
		jwt.WithSecretKey("my-secret-key-123456"),
		jwt.WithIssuer("go-boot-example"),
		jwt.WithExpiresDuration(1*time.Nanosecond), // 极短过期时间
	)
	expiredToken, _, err := expiredJwtUtil.GenerateToken("bob", "test-app")
	if err != nil {
		fmt.Printf("Error generating expired token: %v\n", err)
		return
	}

	// 等待 Token 过期
	time.Sleep(10 * time.Millisecond)

	valid, err = jwtUtil.ValidateToken(expiredToken)
	if err != nil {
		fmt.Printf("Expected error for expired token: %v\n", err)
	} else {
		fmt.Printf("Unexpected result: valid=%v\n", valid)
	}
	fmt.Println()

	fmt.Println("jwt example completed successfully!")
}
