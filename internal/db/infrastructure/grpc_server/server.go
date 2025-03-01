package grpc_server

import (
	"context"
	"github.com/oziev02/checklist-microservices/internal/api/infrastructure/api"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

// Server представляет gRPC-сервер для БД-сервиса
type Server struct {
	grpcServer *grpc.Server
	api        UnimplementedCheckListServiceServer
}

func NewServer() (*Server, error) {
	return &Server{
		grpcServer: grpc.NewServer(),
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
}

func (s *Server) CreateTask(ctx context.Context, req *api.TaskRequest) (*api.TaskResponse, error) {
	// Заглушка: позже реализую с PostgreSQL
	return &api.TaskResponse{
		Id:      "placeholder-task-id",
		Title:   req.Title,
		Content: req.Content,
		Done:    false,
		UserId:  req.UserId,
	}, nil
}

func (s *Server) ListTasks(ctx context.Context, req *api.ListTasksRequest) (*api.ListTasksResponse, error) {
	// Заглушка
	return &api.ListTasksResponse{
		Tasks: []*api.TaskResponse{
			{
				Id:      "placeholder-task-id",
				Title:   "Sample Task",
				Content: "Sample Content",
				Done:    false,
				UserId:  req.UserId,
			},
		},
	}, nil
}

func (s *Server) DeleteTask(ctx context.Context, req *api.TaskIDRequest) (*api.Empty, error) {
	// Заглушка
	return &api.Empty{}, nil
}

func (s *Server) MarkTaskDone(ctx context.Context, req *api.TaskIDRequest) (*api.TaskResponse, error) {
	// Заглушка
	return &api.TaskResponse{
		Id:      req.Id,
		Title:   "Sample Task",
		Content: "Sample Content",
		Done:    true,
		UserId:  "placeholder-user-id",
	}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *api.UserRequest) (*api.UserResponse, error) {
	// Заглушка
	return &api.UserResponse{
		Id:           "placeholder-user-id",
		Email:        req.Email,
		Password:     req.Password,
		Avatar:       "",
		Description:  "",
		Socials:      make(map[string]string),
		TwofaEnabled: false,
		TwofaSecret:  "",
	}, nil
}

func (s *Server) UpdateProfile(ctx context.Context, req *api.UpdateProfileRequest) (*api.UserResponse, error) {
	// Заглушка
	return &api.UserResponse{
		Id:           req.UserId,
		Email:        "placeholder-email@example.com",
		Password:     "placeholder-password",
		Avatar:       req.Avatar,
		Description:  req.Description,
		Socials:      req.Socials,
		TwofaEnabled: false,
		TwofaSecret:  "",
	}, nil
}

func (s *Server) GetUserByEmail(ctx context.Context, req *api.EmailRequest) (*api.UserResponse, error) {
	// Заглушка
	return &api.UserResponse{
		Id:           "placeholder-user-id",
		Email:        req.Email,
		Password:     "placeholder-password",
		Avatar:       "",
		Description:  "",
		Socials:      make(map[string]string),
		TwofaEnabled: false,
		TwofaSecret:  "",
	}, nil
}
