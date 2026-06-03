package comment

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ChristianTertius/devbercerita/internal/dto"
	commentService "ChristianTertius/devbercerita/internal/service/comment"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func TestHandler_CreateComment(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		service     commentService.CommentService
		wantStatus  int
		wantMessage string
	}{
		{
			name: "success",
			body: `{"post_id":2,"content":"hello"}`,
			service: &mockCommentService{
				createFn: func(ctx context.Context, req *dto.StoreCommentRequest, userID int64) (int, error) {
					if req.PostID != 2 || req.Content != "hello" || userID != 7 {
						t.Errorf("unexpected args: %v %v %d", req, req.Content, userID)
					}
					return http.StatusOK, nil
				},
			},
			wantStatus:  http.StatusOK,
			wantMessage: "successfully",
		},
		{
			name:        "validation error",
			body:        `{}`,
			service:     &mockCommentService{},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "StoreCommentRequest",
		},
		{
			name: "service failure",
			body: `{"post_id":5,"content":"boom"}`,
			service: &mockCommentService{
				createFn: func(ctx context.Context, req *dto.StoreCommentRequest, userID int64) (int, error) {
					return http.StatusNotFound, errPostNotFound
				},
			},
			wantStatus:  http.StatusNotFound,
			wantMessage: "post not found",
		},
	}

	header := validator.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			h := NewHandler(r, header, tt.service)
			req := httptest.NewRequest(http.MethodPost, "/comments", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resp)
			ctx.Request = req
			ctx.Set("userID", int64(7))
			h.CreateComment(ctx)

			if resp.Code != tt.wantStatus {
				t.Fatalf("expected status %d got %d", tt.wantStatus, resp.Code)
			}

			var payload map[string]string
			err := json.Unmarshal(resp.Body.Bytes(), &payload)
			if err != nil {
				t.Fatalf("invalid response json: %v", err)
			}

			msg, _ := payload["message"]
			if tt.wantMessage != "" && !strings.Contains(msg, tt.wantMessage) {
				t.Fatalf("expected message to contain %q, got %q", tt.wantMessage, msg)
			}
		})
	}
}

type mockCommentService struct {
	createFn func(ctx context.Context, req *dto.StoreCommentRequest, userID int64) (int, error)
}

func (m *mockCommentService) CreateComment(ctx context.Context, req *dto.StoreCommentRequest, userID int64) (int, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req, userID)
	}
	return http.StatusOK, nil
}

var errPostNotFound = errors.New("post not found")

func (m *mockCommentService) LikeOrUnlikeComment(ctx context.Context, commentID, userID int64) (int, error) {
	return http.StatusOK, nil
}

var _ commentService.CommentService = (*mockCommentService)(nil)
