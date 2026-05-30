package user

import (
	"ChristianTertius/devbercerita/internal/dto"
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *userService) Register(ctx context.Context, req *dto.RegisterRequest) (int64, int, error) {
	// check if user already exists
	userExists, err := s.userRepo.GetUserByEmailOrUsername(ctx, req.Email, req.Username)

	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	if userExists != nil {
		return 0, http.StatusBadRequest, errors.New("user already exists!")
	}

	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	// create user
	now := time.Now()
	userModel := &model.UserModel{
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(passwordHash),
		CreatedAt: now,
		UpdatedAt: now,
	}

	userID, err := s.userRepo.CreateUser(ctx, userModel)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return userID, http.StatusOK, nil
}
