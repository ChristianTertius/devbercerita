package user

import (
	"ChristianTertius/devbercerita/internal/dto"
	"ChristianTertius/devbercerita/internal/model"
	"ChristianTertius/devbercerita/pkg/jwt"
	"ChristianTertius/devbercerita/pkg/refreshtoken"
	"context"
	"errors"
	"net/http"
	"time"
)

func (s *userService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest, userID int64) (string, string, int, error) {
	// check user if exists
	userExist, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	if userExist == nil {
		return "", "", http.StatusNotFound, errors.New("user not found")
	}

	// get refresh token by user id
	refreshTokenExist, err := s.userRepo.GetRefreshToken(ctx, userID, time.Now())
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	if refreshTokenExist == nil {
		return "", "", http.StatusUnauthorized, errors.New("refresh token was expired!")
	}

	// check refresh token is match with request body
	if req.RefreshToken != refreshTokenExist.RefreshToken {
		return "", "", http.StatusUnauthorized, errors.New("refresh token not found!")
	}

	// generate new token
	token, err := jwt.CreateToken(userID, userExist.Username, s.cfg.SecretJwt)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	// delete old refresh token & generate new refresh token
	err = s.userRepo.DeleteRefreshTokenByID(ctx, userID)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	refreshToken, err := refreshtoken.GenerateRefreshToken()
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	now := time.Now()
	s.userRepo.StoreRefreshToken(ctx, &model.RefreshTokenModel{
		UserID:       userID,
		RefreshToken: refreshToken,
		CreatedAt:    now,
		UpdatedAt:    now,
		ExpiredAt:    time.Now().Add(7 * 24 * time.Hour),
	})

	return token, refreshToken, http.StatusOK, nil
}
