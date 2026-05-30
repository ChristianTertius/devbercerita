package user

import (
	"ChristianTertius/devbercerita/internal/config"
	"ChristianTertius/devbercerita/internal/dto"
	"ChristianTertius/devbercerita/internal/repository/user"
	"context"
)

type UserService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (int64, int, error)
}

type userService struct {
	cfg      *config.Config
	userRepo user.UserRepository
}

func NewService(cfg *config.Config, userRepo user.UserRepository) UserService {
	return &userService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}
