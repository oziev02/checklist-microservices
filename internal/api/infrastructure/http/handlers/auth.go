package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oziev02/checklist-microservices/internal/api/domain/entities"
	"github.com/oziev02/checklist-microservices/internal/api/infrastructure/grpc_client"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
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

// authMiddleware проверяет JWT-токен
func (h *AuthHandler) authMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return
	}

	// Проверяем, что заголовок начинается с "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		c.JSON(401, gin.H{"error": "Invalid Authorization header format"})
		c.Abort()
		return
	}

	tokenString := authHeader[7:]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Извлекаем claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	// Сохраняем user_id в контексте для использования в других обработчиках
	userID, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid user_id in token"})
		c.Abort()
		return
	}

	c.Set("user_id", userID)
	c.Next()
}

func (h *AuthHandler) registerHandler(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
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

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(userResp.Password), []byte(req.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	// Генерируем JWT-токен
	accessToken, err := generateToken(userResp.Id, false)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Генерируем Refresh-токен
	refreshToken, err := generateToken(userResp.Id, true)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate refresh token"})
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

func generateToken(userID string, isRefresh bool) (string, error) {
	var secret string
	var expiration time.Duration

	if isRefresh {
		secret = os.Getenv("REFRESH_TOKEN_SECRET")
		expiration = 7 * 24 * time.Hour // 7 дней для Refresh-токена
	} else {
		secret = os.Getenv("JWT_SECRET")
		expiration = 15 * time.Minute // 15 минут для Access-токена
	}

	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET or REFRESH_TOKEN_SECRET not set in .env")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(expiration).Unix(),
	})

	return token.SignedString([]byte(secret))
}
