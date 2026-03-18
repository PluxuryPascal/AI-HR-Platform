package auth_grpc

import (
	"context"

	authv1 "backend/internal/proto/auth/v1"
	"backend/internal/repo"

	"go.uber.org/zap"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	log      *zap.Logger
	userRepo repo.UserRepository
}

func NewAuthHandler(log *zap.Logger, userRepo repo.UserRepository) *AuthHandler {
	return &AuthHandler{
		log:      log,
		userRepo: userRepo,
	}
}

func (h *AuthHandler) GetUsers(ctx context.Context, req *authv1.GetUsersRequest) (*authv1.GetUsersResponse, error) {
	if len(req.UserIds) == 0 {
		return &authv1.GetUsersResponse{
			Users: make(map[string]*authv1.UserObj),
		}, nil
	}

	users, err := h.userRepo.GetUsers(ctx, req.UserIds)
	if err != nil {
		h.log.Error("failed to get users", zap.Error(err))
		return nil, err
	}

	res := &authv1.GetUsersResponse{
		Users: make(map[string]*authv1.UserObj),
	}
	for _, u := range users {
		res.Users[u.ID] = &authv1.UserObj{
			Id:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
		}
	}

	return res, nil
}
