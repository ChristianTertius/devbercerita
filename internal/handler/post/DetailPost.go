package post

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DetailPost returns a single post by ID, including comments.
// @Summary Get post detail
// @Tags Posts
// @Produce json
// @Param post_id path int true "Post ID"
// @Success 200 {object} dto.DetailPostResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 404 {object} dto.MessageResponse
// @Failure 500 {object} dto.MessageResponse
// @Router /posts/{post_id}/detail [get]
func (h *Handler) DetailPost(c *gin.Context) {
	ctx := c.Request.Context()
	postIDParam := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	result, statusCode, err := h.postService.DetailPost(ctx, postID)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(statusCode, result)
}
