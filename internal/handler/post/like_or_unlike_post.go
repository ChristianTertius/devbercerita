package post

import (
	"ChristianTertius/devbercerita/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LikeOrUnlikePost toggles the like status for the requesting user's post.
// @Summary Like or unlike a post
// @Tags Posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.LikeOrUnlikePostRequest true "Post like payload"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 401 {object} dto.MessageResponse
// @Failure 500 {object} dto.MessageResponse
// @Router /posts/action [post]
func (h *Handler) LikeOrUnlikePost(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req dto.LikeOrUnlikePostRequest
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
	statusCode, err := h.postService.LikeOrUnlikePost(ctx, req.PostID, userID)

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
