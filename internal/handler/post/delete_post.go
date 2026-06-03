package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeletePost removes a post owned by the authenticated user.
// @Summary Delete a post
// @Description Delete a post that belongs to the calling user.
// @Tags Posts
// @Security BearerAuth
// @Produce json
// @Param post_id path int true "Post ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 401 {object} dto.MessageResponse
// @Failure 404 {object} dto.MessageResponse
// @Failure 500 {object} dto.MessageResponse
// @Router /posts/{post_id}/delete [delete]
func (h *Handler) DeletePost(c *gin.Context) {
	var (
		ctx       = c.Request.Context()
		userID    = c.GetInt64("userID")
		postIDStr = c.Param("post_id")
	)

	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	statusCode, err := h.postService.DeletePost(ctx, postID, userID)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(statusCode, gin.H{"message": "successfully delete post!"})
}
