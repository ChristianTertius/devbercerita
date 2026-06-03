package user

import (
	"ChristianTertius/devbercerita/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login authenticates a user and returns access and refresh tokens.
// @Summary Login to an existing account
// @Description Verify email and password to issue JWT access and refresh tokens.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login payload"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 404 {object} dto.MessageResponse
// @Failure 500 {object} dto.MessageResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req dto.LoginRequest
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

	token, refreshToken, statusCode, err := h.userService.Login(ctx, &req)

	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(statusCode, dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}
