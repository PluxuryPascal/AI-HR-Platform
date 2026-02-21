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
	return nil, nil
}

func (r *userRepo) Register(ctx context.Context, user *domain.UserRegister) (*domain.User, error) {
	return nil, nil
}
