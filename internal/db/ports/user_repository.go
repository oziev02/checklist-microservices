package ports

import (
	"context"
	"github.com/oziev02/checklist-microservices/internal/db/domain/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, password string) (*entities.User, error)
	UpdateProfile(ctx context.Context, userID, avatar, description string, socials map[string]string) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
}
