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

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	search := c.DefaultQuery("search", "")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	// whitelist validasi
	if sortBy != "created_at" && sortBy != "like_count" {
		sortBy = "created_at"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	param := dto.GetAllPostRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
		SortBy: sortBy,
		Order:  order,
	}

	result, statusCode, err := h.postService.GetAllPost(ctx, &param)
	if err != nil {
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}
	c.JSON(statusCode, result)
}
