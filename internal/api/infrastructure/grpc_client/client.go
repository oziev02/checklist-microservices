package grpc_client

import (
	"context"
	"github.com/oziev02/checklist-microservices/internal/api/infrastructure/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

// Client представляет gRPC-клиент для взаимодействия с БД-сервисом
type Client struct {
	conn    *grpc.ClientConn
	service api.ChecklistServiceClient
}

func NewClient() (*Client, error) {
	// Получаем хост и порт БД-сервиса из переменных окружения
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	if host == "" || port == "" {
		log.Fatalf("DB_HOST and DB_PORT must be set in .env")
	}

	// Подключаемся к БД-сервису
	conn, err := grpc.Dial(host+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to DB service: %v", err)
	}

	// Создаём gRPC-клиент
	service := api.NewChecklistServiceClient(conn)

	return &Client{
		conn:    conn,
		service: service,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateTask(ctx context.Context, title, content, userID string) (*api.TaskResponse, error) {
	req := &api.TaskRequest{
		Title:   title,
		Content: content,
		UserId:  userID,
	}
	return c.service.CreateTask(ctx, req)
}

func (c *Client) ListTasks(ctx context.Context, userID string) (*api.ListTasksResponse, error) {
	req := &api.ListTasksRequest{
		UserId: userID,
	}
	return c.service.ListTasks(ctx, req)
}

func (c *Client) DeleteTask(ctx context.Context, taskID string) error {
	req := &api.TaskIDRequest{
		Id: taskID,
	}
	_, err := c.service.DeleteTask(ctx, req)
	return err
}

func (c *Client) MarkTaskDone(ctx context.Context, taskID string) (*api.TaskResponse, error) {
	req := &api.TaskIDRequest{
		Id: taskID,
	}
	return c.service.MarkTaskDone(ctx, req)
}

func (c *Client) CreateUser(ctx context.Context, email, password string) (*api.UserResponse, error) {
	req := &api.UserRequest{
		Email:    email,
		Password: password,
	}
	return c.service.CreateUser(ctx, req)
}

func (c *Client) UpdateProfile(ctx context.Context, userID, avatar, description string, socials map[string]string) (*api.UserResponse, error) {
	req := &api.UpdateProfileRequest{
		UserId:      userID,
		Avatar:      avatar,
		Description: description,
		Socials:     socials,
	}
	return c.service.UpdateProfile(ctx, req)
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*api.UserResponse, error) {
	req := &api.EmailRequest{
		Email: email,
	}
	return c.service.GetUserByEmail(ctx, req)
}
