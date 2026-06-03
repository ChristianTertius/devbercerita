package post

import (
	"ChristianTertius/devbercerita/internal/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdatePost modifies an existing post that belongs to the authenticated user.
// @Summary Update a post
// @Description Update the title or content of an owned post.
// @Tags Posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param request body dto.CreateOrUpdatePostRequest true "Updated post payload"
// @Success 200 {object} dto.CreateOrUpdatePostResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 401 {object} dto.MessageResponse
// @Failure 404 {object} dto.MessageResponse
// @Failure 500 {object} dto.MessageResponse
// @Router /posts/{post_id}/update [put]
func (h *Handler) UpdatePost(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req dto.CreateOrUpdatePostRequest
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
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	statusCode, err := h.postService.UpdatePost(ctx, &req, postID, userID)

	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(statusCode, dto.CreateOrUpdatePostResponse{
		ID: postID,
	})
}
