package repo

import (
	"backend/internal/db"
	"backend/internal/domain"
	"context"
)

type UserRepository interface {
	Login(ctx context.Context, login string) (*domain.UserLogin, error)
	Register(ctx context.Context, user *domain.UserRegister) (*domain.User, error)
}

type userRepo struct {
	dbClient *db.Client
}

func NewUserRepo(dbClient *db.Client) UserRepository {
	return &userRepo{dbClient: dbClient}
}

func (r *userRepo) Login(ctx context.Context, login string) (*domain.UserLogin, error) {
	return &domain.UserLogin{
		ID:         "mock-id-123",
		Password:   "$argon2id$v=19$m=65536,t=1,p=4$4wDt2Cx/ryTr0/THdujebA$hCLlQYObN74RQSKB9ZXImegzFOV+UdWdHvi+9iDiV0g", // argon2 hash of 'password123'
		GroupAlias: "admin",
	}, nil
}

func (r *userRepo) Register(ctx context.Context, user *domain.UserRegister) (*domain.User, error) {
	return &domain.User{
		ID:         "mock-id-123",
		GroupAlias: "admin",
	}, nil
}
