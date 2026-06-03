package post

import (
	"ChristianTertius/devbercerita/internal/dto"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllPost returns paginated post summaries.
// @Summary List posts
// @Description Returns paginated posts with default limit 10 and page 1.
// @Tags Posts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.GetAllPostResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 500 {object} dto.MessageResponse
// @Router /posts [get]
func (h *Handler) GetAllPost(c *gin.Context) {
	ctx := c.Request.Context()
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.ParseInt(pageStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	param := dto.GetAllPostRequest{
		Page:  page,
		Limit: limit,
	}

	result, statusCode, err := h.postService.GetAllPost(ctx, &param)
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(statusCode, result)
}
