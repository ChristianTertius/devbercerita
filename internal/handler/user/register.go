package user

import (
	"ChristianTertius/devbercerita/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register registers a new user account via email, username, and password.
// @Summary Register a new user
// @Description Create an account with email, username, and password credentials.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration payload"
// @Success 200 {object} dto.RegisterResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 500 {object} dto.MessageResponse
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req dto.RegisterRequest
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

	userID, statusCode, err := h.userService.Register(ctx, &req)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(statusCode, dto.RegisterResponse{
		ID: userID,
	})
}
