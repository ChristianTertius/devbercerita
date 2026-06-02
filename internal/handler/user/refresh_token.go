package user

import (
	"ChristianTertius/devbercerita/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RefreshToken rotates access and refresh tokens after validating the provided refresh token.
// @Summary Refresh session tokens
// @Description Exchange an existing refresh token (via Authorization header) for new tokens.
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token payload"
// @Success 200 {object} dto.RefreshTokenResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 401 {object} dto.MessageResponse
// @Failure 404 {object} dto.MessageResponse
// @Failure 500 {object} dto.MessageResponse
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req dto.RefreshTokenRequest
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	userID := c.GetInt64("userID")
	token, refreshToken, statusCode, err := h.userService.RefreshToken(ctx, &req, userID)

	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(statusCode, dto.RefreshTokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}
