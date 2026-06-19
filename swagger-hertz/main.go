// Package main 演示 go-boot 与 Hertz + Swagger 集成
//
// 本示例展示如何:
//   - 创建 Hertz 服务器并与 go-boot 集成
//   - 集成 Swagger 文档
//   - 配置 Swagger UI
//   - 定义 API 路由和文档注释
//
// 使用方式:
//
//	cd examples/swagger-hertz
//	go run .
//	# 访问 http://localhost:8081/swagger/index.html
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/xudefa/go-boot/core"
	hs "github.com/xudefa/go-boot-hertz/server"
)

// @title           Go-Boot Hertz API
// @version         1.0
// @description     这是一个使用 go-boot + Hertz + Swagger 的示例 API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8081
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Bearer token authentication

func main() {
	fmt.Println("=== Go-Boot Hertz + Swagger Example ===")
	fmt.Println()

	container := core.New()
	fmt.Println("Container created")

	server := hs.NewServer(hs.WithHost("localhost"), hs.WithPort(8081))
	fmt.Println("Hertz server created")

	api := server.Group("/api/v1")
	{
		api.GET("/hello", helloHandler)
		api.GET("/users/:id", getUserHandler)
		api.POST("/users", createUserHandler)
	}

	fmt.Println("Routes registered:")
	fmt.Println("  GET /api/v1/hello")
	fmt.Println("  GET /api/v1/users/:id")
	fmt.Println("  POST /api/v1/users")

	server.GET("/swagger/*any", swaggerHandler)

	err := container.Register("hertzServer", core.Bean(server))
	if err != nil {
		log.Printf("Warning: Failed to register server: %v\n", err)
	}

	fmt.Println("\nStarting server on :8081...")
	fmt.Println("API Documentation: http://localhost:8081/swagger/index.html")
	fmt.Println("Press Ctrl+C to stop")

	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// helloHandler 返回欢迎信息
// @Summary      获取欢迎信息
// @Description  返回简单的欢迎消息
// @Tags         hello
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /hello [get]
func helloHandler(ctx context.Context, c *app.RequestContext) {
	c.JSON(200, utils.H{
		"message": "Hello from go-boot with Hertz + Swagger!",
	})
}

// getUserHandler 获取用户信息
// @Summary      获取用户信息
// @Description  根据 ID 获取用户详细信息
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "用户 ID"
// @Success      200  {object}  UserResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /users/{id} [get]
func getUserHandler(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	c.JSON(200, utils.H{
		"id":    id,
		"name":  "John Doe",
		"email": "john@example.com",
	})
}

// createUserHandler 创建新用户
// @Summary      创建新用户
// @Description  创建一个新用户并返回用户信息
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "用户信息"
// @Success      201   {object}  UserResponse
// @Failure      400   {object}  ErrorResponse
// @Router       /users [post]
func createUserHandler(ctx context.Context, c *app.RequestContext) {
	var req CreateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, utils.H{"error": err.Error()})
		return
	}
	c.JSON(201, utils.H{
		"id":    "1",
		"name":  req.Name,
		"email": req.Email,
	})
}

// swaggerHandler 处理 Swagger UI 请求
func swaggerHandler(ctx context.Context, c *app.RequestContext) {
	requestPath := string(c.Request.URI().Path())

	if requestPath == "/swagger/" || requestPath == "/swagger" {
		c.SetContentType("text/html; charset=utf-8")
		c.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui.css">
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/swagger/doc.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`)
		return
	}

	if requestPath == "/swagger/doc.json" {
		c.SetContentType("application/json; charset=utf-8")
		c.WriteString(`{
  "swagger": "2.0",
  "info": {
    "title": "Go-Boot Hertz API",
    "description": "这是一个使用 go-boot + Hertz + Swagger 的示例 API",
    "version": "1.0",
    "contact": {
      "name": "API Support",
      "email": "support@swagger.io",
      "url": "http://www.swagger.io/support"
    },
    "license": {
      "name": "MIT",
      "url": "https://opensource.org/licenses/MIT"
    }
  },
  "host": "localhost:8081",
  "basePath": "/api/v1",
  "schemes": ["http"],
  "paths": {
    "/hello": {
      "get": {
        "tags": ["hello"],
        "summary": "获取欢迎信息",
        "description": "返回简单的欢迎消息",
        "produces": ["application/json"],
        "responses": {
          "200": {
            "description": "成功",
            "schema": {
              "type": "object",
              "properties": {
                "message": {
                  "type": "string",
                  "example": "Hello from go-boot with Hertz + Swagger!"
                }
              }
            }
          }
        }
      }
    },
    "/users/{id}": {
      "get": {
        "tags": ["users"],
        "summary": "获取用户信息",
        "description": "根据 ID 获取用户详细信息",
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "type": "string",
            "description": "用户 ID"
          }
        ],
        "responses": {
          "200": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/UserResponse"
            }
          },
          "404": {
            "description": "用户不存在",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/users": {
      "post": {
        "tags": ["users"],
        "summary": "创建新用户",
        "description": "创建一个新用户并返回用户信息",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "name": "user",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/CreateUserRequest"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "创建成功",
            "schema": {
              "$ref": "#/definitions/UserResponse"
            }
          },
          "400": {
            "description": "请求参数错误",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "CreateUserRequest": {
      "type": "object",
      "required": ["name", "email"],
      "properties": {
        "name": {
          "type": "string",
          "example": "John Doe"
        },
        "email": {
          "type": "string",
          "format": "email",
          "example": "john@example.com"
        }
      }
    },
    "UserResponse": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "example": "1"
        },
        "name": {
          "type": "string",
          "example": "John Doe"
        },
        "email": {
          "type": "string",
          "example": "john@example.com"
        }
      }
    },
    "ErrorResponse": {
      "type": "object",
      "properties": {
        "error": {
          "type": "string",
          "example": "invalid request"
        }
      }
    }
  },
  "securityDefinitions": {
    "Bearer": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header",
      "description": "Bearer token authentication"
    }
  },
  "security": [
    {
      "Bearer": []
    }
  ]
}`)
		return
	}

	c.SetStatusCode(404)
	c.WriteString("404 Not Found")
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID    string `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request"`
}
