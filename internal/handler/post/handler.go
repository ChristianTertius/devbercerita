package post

import (
	"ChristianTertius/devbercerita/internal/middleware"
	"ChristianTertius/devbercerita/internal/service/post"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	api         *gin.Engine
	validate    *validator.Validate
	postService post.PostService
}

func NewHandler(api *gin.Engine, validate *validator.Validate, postService post.PostService) *Handler {
	return &Handler{
		api:         api,
		validate:    validate,
		postService: postService,
	}
}

func (h *Handler) RouteList(secretKey string) {
	routeAuth := h.api.Group("/posts")
	routeAuth.Use(middleware.AuthMiddleware(secretKey))
	routeAuth.POST("/", h.CreatePost)
	routeAuth.PUT("/:post_id/update", h.UpdatePost)
	routeAuth.DELETE("/:post_id/delete", h.DeletePost)
	routeAuth.POST("/action", h.LikeOrUnlikePost)

	routeWithoutAuth := h.api.Group("/posts")
	routeWithoutAuth.GET("/:post_id/detail", h.DetailPost)
	routeWithoutAuth.GET("/", h.GetAllPost)
}
