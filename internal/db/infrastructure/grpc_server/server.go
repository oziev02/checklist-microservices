package grpc_server

import (
	"context"
	"fmt"
	"github.com/oziev02/checklist-microservices/internal/api/infrastructure/api"
	"github.com/oziev02/checklist-microservices/internal/db/infrastructure/database"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

// Server представляет gRPC-сервер для БД-сервиса
type Server struct {
	grpcServer *grpc.Server
	repo       *database.PostgresRepository
	api.UnimplementedChecklistServiceServer
}

func NewServer() (*Server, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatalf("POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, and POSTGRES_DB must be set in .env")
	}

	repo, err := database.NewPostgresRepository(host, port, user, password, dbname)
	if err != nil {
		log.Fatalf("Failed to create Postgres repository: %v", err)
	}

	return &Server{
		grpcServer: grpc.NewServer(),
		repo:       repo,
	}, nil
}

// Start запускает gRPC-сервер
func (s *Server) Start() error {
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "8081"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Регистрируем сервис
	api.RegisterChecklistServiceServer(s.grpcServer, s)

	// Запускаем сервер
	log.Printf("Starting gRPC server on port %s", port)
	return s.grpcServer.Serve(lis)
}

// Stop останавливает gRPC-сервер
func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
	if err := s.repo.Close(); err != nil {
		log.Printf("Failed to close Postgres repository: %v", err)
	}
}

func (s *Server) CreateTask(ctx context.Context, req *api.TaskRequest) (*api.TaskResponse, error) {
	task, err := s.repo.CreateTask(ctx, req.Title, req.Content, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("Failed to create task: %v", err)
	}
	return &api.TaskResponse{
		Id:      task.ID.String(),
		Title:   task.Title,
		Content: task.Content,
		Done:    task.Done,
		UserId:  task.UserID.String(),
	}, nil
}

func (s *Server) ListTasks(ctx context.Context, req *api.ListTasksRequest) (*api.ListTasksResponse, error) {
	tasks, err := s.repo.ListTasks(ctx, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("Failed to list tasks: %v", err)
	}

	var taskResponses []*api.TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, &api.TaskResponse{
			Id:      task.ID.String(),
			Title:   task.Title,
			Content: task.Content,
			Done:    task.Done,
			UserId:  task.UserID.String(),
		})
	}

	return &api.ListTasksResponse{Tasks: taskResponses}, nil
}

func (s *Server) DeleteTask(ctx context.Context, req *api.TaskIDRequest) (*api.Empty, error) {
	if err := s.repo.DeleteTask(ctx, req.Id); err != nil {
		return nil, fmt.Errorf("Failed to delete task: %v", err)
	}
	return &api.Empty{}, nil
}

// Обрабатывает отметку задачи как выполненной
func (s *Server) MarkTaskDone(ctx context.Context, req *api.TaskIDRequest) (*api.TaskResponse, error) {
	task, err := s.repo.MarkTaskDone(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("Failed to mark task as done: %v", err)
	}
	return &api.TaskResponse{
		Id:      task.ID.String(),
		Title:   task.Title,
		Content: task.Content,
		Done:    task.Done,
		UserId:  task.UserID.String(),
	}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *api.UserRequest) (*api.UserResponse, error) {
	user, err := s.repo.CreateUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user: %v", err)
	}
	return &api.UserResponse{
		Id:           user.ID.String(),
		Email:        user.Email,
		Password:     user.Password,
		Avatar:       user.Avatar,
		Description:  user.Description,
		Socials:      user.Socials,
		TwofaEnabled: user.TwoFAEnabled,
		TwofaSecret:  user.TwoFASecret,
	}, nil
}

func (s *Server) UpdateProfile(ctx context.Context, req *api.UpdateProfileRequest) (*api.UserResponse, error) {
	user, err := s.repo.UpdateProfile(ctx, req.UserId, req.Avatar, req.Description, req.Socials)
	if err != nil {
		return nil, fmt.Errorf("Failed to update profile: %v", err)
	}
	return &api.UserResponse{
		Id:           user.ID.String(),
		Email:        user.Email,
		Password:     user.Password,
		Avatar:       user.Avatar,
		Description:  user.Description,
		Socials:      user.Socials,
		TwofaEnabled: user.TwoFAEnabled,
		TwofaSecret:  user.TwoFASecret,
	}, nil
}

func (s *Server) GetUserByEmail(ctx context.Context, req *api.EmailRequest) (*api.UserResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user by email: %v", err)
	}
	if user == nil {
		return nil, nil // Пользователь не найден
	}
	return &api.UserResponse{
		Id:           user.ID.String(),
		Email:        user.Email,
		Password:     user.Password,
		Avatar:       user.Avatar,
		Description:  user.Description,
		Socials:      user.Socials,
		TwofaEnabled: user.TwoFAEnabled,
		TwofaSecret:  user.TwoFASecret,
	}, nil
}
