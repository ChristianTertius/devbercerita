package post

import (
	"ChristianTertius/devbercerita/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePost saves a new post for the authenticated user.
// @Summary Create a post
// @Description Authenticated user can create a new post with title and content.
// @Tags Posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateOrUpdatePostRequest true "Post payload"
// @Success 200 {object} dto.CreateOrUpdatePostResponse
// @Failure 400 {object} dto.MessageResponse
// @Failure 401 {object} dto.MessageResponse
// @Failure 500 {object} dto.MessageResponse
// @Router /posts [post]
func (h *Handler) CreatePost(c *gin.Context) {
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
	postID, statusCode, err := h.postService.CreatePost(ctx, &req, userID)

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
