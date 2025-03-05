package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/oziev02/checklist-microservices/internal/api/infrastructure/grpc_client"
)

type ProfileHandler struct {
	grpcClient *grpc_client.Client
}

func NewProfileHandler(grpcClient *grpc_client.Client) *ProfileHandler {
	return &ProfileHandler{grpcClient: grpcClient}
}

func (h *ProfileHandler) getProfileHandler(c *gin.Context) {
	// Получаем userID из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "User ID not found in context"})
		return
	}

	// Отправляем запрос в БД-сервис через gRPC
	userResp, err := h.grpcClient.GetUserByEmail(context.Background(), "placeholder-email@example.com") // Заменим позже на поиск по userID
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user profile"})
		return
	}
	if userResp == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{
		"id":          userResp.Id,
		"email":       userResp.Email,
		"avatar":      userResp.Avatar,
		"description": userResp.Description,
		"socials":     userResp.Socials,
	})
}

func (h *ProfileHandler) updateProfileHandler(c *gin.Context) {
	var req struct {
		Avatar      string            `json:"avatar"`
		Description string            `json:"description"`
		Socials     map[string]string `json:"socials"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Получаем userID из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "User ID not found in context"})
		return
	}

	// Отправляем запрос в БД-сервис через gRPC
	userResp, err := h.grpcClient.UpdateProfile(context.Background(), userID.(string), req.Avatar, req.Description, req.Socials)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(200, gin.H{
		"id":          userResp.Id,
		"email":       userResp.Email,
		"avatar":      userResp.Avatar,
		"description": userResp.Description,
		"socials":     userResp.Socials,
	})
}
