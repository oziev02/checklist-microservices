package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/oziev02/checklist-microservices/internal/api/domain/entities"
	"github.com/oziev02/checklist-microservices/internal/api/infrastructure/grpc_client"
)

type AuthHandler struct {
	grpcClient *grpc_client.Client
}

func NewAuthHandler(grpcClient *grpc_client.Client) *AuthHandler {
	return &AuthHandler{grpcClient: grpcClient}
}

// Регистрирует все маршруты API
func RegisterRoutes(r *gin.Engine, grpcClient *grpc_client.Client) {
	authHandler := NewAuthHandler(grpcClient)
	taskHandler := NewTaskHandler(grpcClient)
	profileHandler := NewProfileHandler(grpcClient)

	// Аутентификация
	r.POST("/register", authHandler.registerHandler)
	r.POST("/login", authHandler.loginHandler)
	r.POST("/2fa/setup", authHandler.setup2FAHandler)
	r.POST("/2fa/verify", authHandler.verify2FAHandler)

	auth := r.Group("/", authHandler.authMiddleware)
	{
		// Tasks
		auth.POST("/create", taskHandler.createTaskHandler)
		auth.GET("/list", taskHandler.listTasksHandler)
		auth.DELETE("/delete", taskHandler.deleteTaskHandler)
		auth.PUT("/done", taskHandler.markTaskDoneHandler)

		// Profile
		auth.GET("/profile", profileHandler.getProfileHandler)
		auth.PUT("/profile", profileHandler.updateProfileHandler)
	}
}

// authMiddleware проверяет JWT-токен (пока это заглушка)
func (h *AuthHandler) authMiddleware(c *gin.Context) {
	// Пока просто пропускаем запрос
	c.Next()
}

func (h *AuthHandler) registerHandler(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	// Отправляем запрос в БД-сервис через gRPC
	resp, err := h.grpcClient.CreateUser(context.Background(), user.Email, user.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(200, gin.H{
		"id":    resp.Id,
		"email": resp.Email,
	})
}

func (h *AuthHandler) loginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Отправляем запрос в БД-сервис через gRPC
	userResp, err := h.grpcClient.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user"})
		return
	}
	if userResp == nil || userResp.Password != req.Password {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  "placeholder-access-token",
		"refresh_token": "placeholder-refresh-token",
	})
}

// настраивает 2FA для пользователя
func (h *AuthHandler) setup2FAHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "2FA setup endpoint placeholder"})
}

// проверяет 2FA-код
func (h *AuthHandler) verify2FAHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "2FA verify endpoint placeholder"})
}
