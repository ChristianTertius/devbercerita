package comment

import (
	"ChristianTertius/devbercerita/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LikeOrUnlikeComment toggles a like on an existing comment.
// @Summary Like or unlike a comment
// @Tags Comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.LikeOrUnlikeCommentRequest true "Comment like payload"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 401 {object} dto.MessageResponse
// @Failure 500 {object} dto.MessageResponse
// @Router /comments/action [post]
func (h *Handler) LikeOrUnlikeComment(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req dto.LikeOrUnlikeCommentRequest
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
	statusCode, err := h.commentService.LikeOrUnlikeComment(ctx, req.CommentID, userID)

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
