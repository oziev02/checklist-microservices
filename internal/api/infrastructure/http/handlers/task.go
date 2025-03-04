package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/oziev02/checklist-microservices/internal/api/infrastructure/grpc_client"
)

type TaskHandler struct {
	grpcClient *grpc_client.Client
}

func NewTaskHandler(grpcClient *grpc_client.Client) *TaskHandler {
	return &TaskHandler{grpcClient: grpcClient}
}

func (h *TaskHandler) createTaskHandler(c *gin.Context) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Имитация userID (позже получим из JWT)
	userID := "placeholder-user-id"

	// Отправляем запрос в БД-сервис через gRPC
	resp, err := h.grpcClient.CreateTask(context.Background(), req.Title, req.Content, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(200, gin.H{
		"id":      resp.Id,
		"title":   resp.Title,
		"content": resp.Content,
		"done":    resp.Done,
		"user_id": resp.UserId,
	})
}

func (h *TaskHandler) listTasksHandler(c *gin.Context) {
	// Имитация userID (позже получим из JWT)
	userID := "placeholder-user-id"

	// Отправляем запрос в БД-сервис через gRPC
	resp, err := h.grpcClient.ListTasks(context.Background(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list tasks"})
		return
	}

	c.JSON(200, resp.Tasks)
}

func (h *TaskHandler) deleteTaskHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Task ID is required"})
		return
	}

	// Отправляем запрос в БД-сервис через gRPC
	if err := h.grpcClient.DeleteTask(context.Background(), id); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete task"})
		return
	}

	c.Status(204)
}

func (h *TaskHandler) markTaskDoneHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Task ID is required"})
		return
	}

	// Отправляем запрос в БД-сервис через gRPC
	resp, err := h.grpcClient.MarkTaskDone(context.Background(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to mark task as done"})
		return
	}

	c.JSON(200, gin.H{
		"id":      resp.Id,
		"title":   resp.Title,
		"content": resp.Content,
		"done":    resp.Done,
		"user_id": resp.UserId,
	})
}
