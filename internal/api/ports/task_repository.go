package ports

import (
	"context"
	"github.com/oziev02/checklist-microservices/internal/api/domain/entities"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, title, content, userID string) (*entities.Task, error)
	ListTasks(ctx context.Context, userID string) ([]*entities.Task, error)
	DeleteTask(ctx context.Context, taskID string) error
	MarkTaskDone(ctx context.Context, taskID string) (*entities.Task, error)
}
