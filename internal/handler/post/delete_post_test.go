package post

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ChristianTertius/devbercerita/internal/dto"
	"ChristianTertius/devbercerita/internal/service/post"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func TestHandler_DeletePost(t *testing.T) {
	tests := []struct {
		name        string
		param       string
		userID      int64
		service     post.PostService
		wantStatus  int
		wantMessage string
		expectCall  bool
	}{
		{
			name:   "success",
			param:  "3",
			userID: 5,
			service: &mockPostService{
				deleteFn: func(ctx context.Context, postID, userID int64) (int, error) {
					if postID != 3 || userID != 5 {
						t.Fatalf("unexpected args: %d %d", postID, userID)
					}
					return http.StatusOK, nil
				},
			},
			wantStatus:  http.StatusOK,
			wantMessage: "successfully delete post",
			expectCall:  true,
		},
		{
			name:        "invalid id",
			param:       "abc",
			userID:      2,
			service:     &mockPostService{},
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "invalid syntax",
			expectCall:  false,
		},
		{
			name:   "service error",
			param:  "11",
			userID: 2,
			service: &mockPostService{
				deleteFn: func(ctx context.Context, postID, userID int64) (int, error) {
					return http.StatusNotFound, errPostNotFound
				},
			},
			wantStatus:  http.StatusNotFound,
			wantMessage: "post not found",
			expectCall:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			h := NewHandler(r, validator.New(), tt.service)
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodDelete, "/posts/"+tt.param+"/delete", nil)
			ctx.Params = gin.Params{{Key: "post_id", Value: tt.param}}
			ctx.Set("userID", tt.userID)
			h.DeletePost(ctx)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected %d got %d", tt.wantStatus, rec.Code)
			}

			var body map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("invalid json: %v", err)
			}

			msg, _ := body["message"]
			if tt.wantMessage != "" && !contains(msg, tt.wantMessage) {
				t.Fatalf("expected message %q contains %q", msg, tt.wantMessage)
			}

			if tt.expectCall {
				if tt.service.(*mockPostService).called == 0 {
					t.Fatalf("expected service call")
				}
			} else {
				if tt.service.(*mockPostService).called > 0 {
					t.Fatalf("service should not be called")
				}
			}
		})
	}
}

func contains(text, sub string) bool {
	return strings.Contains(text, sub)
}

type mockPostService struct {
	called   int
	deleteFn func(ctx context.Context, postID, userID int64) (int, error)
}

func (m *mockPostService) CreatePost(ctx context.Context, req *dto.CreateOrUpdatePostRequest, userID int64) (int64, int, error) {
	return 0, 0, nil
}

func (m *mockPostService) UpdatePost(ctx context.Context, req *dto.CreateOrUpdatePostRequest, postID, userID int64) (int, error) {
	return 0, nil
}

func (m *mockPostService) DeletePost(ctx context.Context, postID, userID int64) (int, error) {
	m.called++
	if m.deleteFn != nil {
		return m.deleteFn(ctx, postID, userID)
	}
	return http.StatusOK, nil
}

func (m *mockPostService) LikeOrUnlikePost(ctx context.Context, postID, userID int64) (int, error) {
	return 0, nil
}

func (m *mockPostService) DetailPost(ctx context.Context, postID int64) (*dto.DetailPostResponse, int, error) {
	return nil, 0, nil
}

func (m *mockPostService) GetAllPost(ctx context.Context, param *dto.GetAllPostRequest) (*dto.GetAllPostResponse, int, error) {
	return nil, 0, nil
}

var errPostNotFound = errors.New("post not found")

var _ post.PostService = (*mockPostService)(nil)
