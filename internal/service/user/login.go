package user

import (
	"ChristianTertius/devbercerita/internal/dto"
	"ChristianTertius/devbercerita/internal/model"
	"ChristianTertius/devbercerita/pkg/jwt"
	"ChristianTertius/devbercerita/pkg/refreshtoken"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *userService) Login(ctx context.Context, req *dto.LoginRequest) (string, string, int, error) {
	// cek user if exists
	userExist, err := s.userRepo.GetUserByEmailOrUsername(ctx, req.Email, "")
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	if userExist == nil {
		return "", "", http.StatusNotFound, errors.New("wrong email or password")
	}

	now := time.Now()

	err = bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(req.Password))
	if err != nil {
		if isBycryptHash(userExist.Password) {
			return "", "", http.StatusNotFound, errors.New("wrong email or password")
		}
		if userExist.Password != req.Password {
			return "", "", http.StatusNotFound, errors.New("wrong email or password")
		}
		hashedPassword, hasErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if hasErr != nil {
			return "", "", http.StatusInternalServerError, hasErr
		}

		userExist.Password = string(hashedPassword)
	}

	// generate access token
	token, err := jwt.CreateToken(userExist.ID, userExist.Username, s.cfg.SecretJwt)
	if err != nil {
		return "", "", http.StatusBadRequest, err
	}

	// get refresh token if exists
	refreshTokenExist, err := s.userRepo.GetRefreshToken(ctx, userExist.ID, now)

	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	if refreshTokenExist != nil {
		return token, refreshTokenExist.RefreshToken, http.StatusOK, nil
	}

	// generate & restore refresh token
	refreshToken, err := refreshtoken.GenerateRefreshToken()
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	err = s.userRepo.StoreRefreshToken(ctx, &model.RefreshTokenModel{
		UserID:       userExist.ID,
		RefreshToken: refreshToken,
		CreatedAt:    now,
		UpdatedAt:    now,
		ExpiredAt:    time.Now().Add(7 * 24 * time.Hour),
	})

	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	return token, refreshToken, http.StatusOK, nil
}

// cek apakah passwordnya itu bcrypt or not
func isBycryptHash(value string) bool {
	if len(value) < 4 {
		return false

	}
	return strings.HasPrefix(value, "$2a$") || strings.HasPrefix(value, "$2b$") || strings.HasPrefix(value, "$2y$")
}
