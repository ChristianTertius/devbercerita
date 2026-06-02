package comment

import (
	"ChristianTertius/devbercerita/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateComment creates a comment tied to a post for the current user.
// @Summary Create a comment
// @Tags Comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.StoreCommentRequest true "Comment payload"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 401 {object} dto.MessageResponse
// @Failure 500 {object} dto.MessageResponse
// @Router /comments [post]
func (h *Handler) CreateComment(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req dto.StoreCommentRequest
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
	statusCode, err := h.commentService.CreateComment(ctx, &req, userID)

	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(statusCode, gin.H{
		"message": "successfully",
	})
}
