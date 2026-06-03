package post

import (
	"context"
	"errors"
	"testing"
	"time"

	"ChristianTertius/devbercerita/internal/dto"
	"ChristianTertius/devbercerita/internal/model"
	"net/http"
)

func TestPostService_DeletePost(t *testing.T) {
	tests := []struct {
		name        string
		postID      int64
		userID      int64
		postRepo    *testPostRepo
		commentRepo *testCommentRepo
		wantStatus  int
		wantErr     bool
	}{
		{
			name:   "success",
			postID: 1,
			userID: 1,
			postRepo: &testPostRepo{
				getFn: func(ctx context.Context, postID int64) (*model.PostWithUserModel, error) {
					return &model.PostWithUserModel{ID: postID, UserID: 1}, nil
				},
				softDeleteFn: func(ctx context.Context, postID int64, now time.Time) error {
					return nil
				},
			},
			commentRepo: &testCommentRepo{deleteFn: func(ctx context.Context, postID int64) error {
				return nil
			}},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:   "post not found",
			postID: 2,
			userID: 1,
			postRepo: &testPostRepo{
				getFn: func(ctx context.Context, postID int64) (*model.PostWithUserModel, error) {
					return nil, nil
				},
			},
			commentRepo: &testCommentRepo{},
			wantStatus:  http.StatusNotFound,
			wantErr:     true,
		},
		{
			name:   "user mismatch",
			postID: 3,
			userID: 99,
			postRepo: &testPostRepo{
				getFn: func(ctx context.Context, postID int64) (*model.PostWithUserModel, error) {
					return &model.PostWithUserModel{ID: postID, UserID: 1}, nil
				},
			},
			commentRepo: &testCommentRepo{},
			wantStatus:  http.StatusNotFound,
			wantErr:     true,
		},
		{
			name:   "delete comments failure",
			postID: 4,
			userID: 1,
			postRepo: &testPostRepo{
				getFn: func(ctx context.Context, postID int64) (*model.PostWithUserModel, error) {
					return &model.PostWithUserModel{ID: postID, UserID: 1}, nil
				},
				softDeleteFn: func(ctx context.Context, postID int64, now time.Time) error {
					return nil
				},
			},
			commentRepo: &testCommentRepo{deleteFn: func(ctx context.Context, postID int64) error {
				return errors.New("boom")
			}},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:   "soft delete failure",
			postID: 5,
			userID: 1,
			postRepo: &testPostRepo{
				getFn: func(ctx context.Context, postID int64) (*model.PostWithUserModel, error) {
					return &model.PostWithUserModel{ID: postID, UserID: 1}, nil
				},
				softDeleteFn: func(ctx context.Context, postID int64, now time.Time) error {
					return errors.New("boom")
				},
			},
			commentRepo: &testCommentRepo{deleteFn: func(ctx context.Context, postID int64) error {
				return nil
			}},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(nil, tt.postRepo, tt.commentRepo)
			gotStatus, err := s.DeletePost(context.Background(), tt.postID, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error state: %v", err)
			}

			if gotStatus != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, gotStatus)
			}
		})
	}
}

type testPostRepo struct {
	getFn        func(ctx context.Context, postID int64) (*model.PostWithUserModel, error)
	softDeleteFn func(ctx context.Context, postID int64, now time.Time) error
}

func (r *testPostRepo) StorePost(ctx context.Context, model *model.PostModel) (int64, error) {
	return 0, nil
}

func (r *testPostRepo) GetPostById(ctx context.Context, postID int64) (*model.PostWithUserModel, error) {
	if r.getFn == nil {
		return nil, nil
	}
	return r.getFn(ctx, postID)
}

func (r *testPostRepo) UpdatePost(ctx context.Context, model *model.PostModel, postID int64) error {
	return nil
}

func (r *testPostRepo) SoftDelete(ctx context.Context, postID int64, now time.Time) error {
	if r.softDeleteFn == nil {
		return nil
	}
	return r.softDeleteFn(ctx, postID, now)
}

func (r *testPostRepo) IsUserAlreadyLikePost(ctx context.Context, postID, userID int64) (bool, error) {
	return false, nil
}

func (r *testPostRepo) DeleteLikePost(ctx context.Context, postID, userID int64) error {
	return nil
}

func (r *testPostRepo) StoreLikePost(ctx context.Context, model *model.PostLikeModel) error {
	return nil
}

func (r *testPostRepo) TotalPost(ctx context.Context) (int64, error) {
	return 0, nil
}

func (r *testPostRepo) GetAllPost(ctx context.Context, param *dto.GetAllPostRequest, offset int) ([]model.PostWithUserModel, error) {
	return nil, nil
}

type testCommentRepo struct {
	deleteFn func(ctx context.Context, postID int64) error
}

func (r *testCommentRepo) StoreComment(ctx context.Context, model *model.CommentModel) error {
	return nil
}
func (r *testCommentRepo) DetailComment(ctx context.Context, commentID int64) (*model.CommentModel, error) {
	return nil, nil
}
func (r *testCommentRepo) IsUserAlreadyLikeComment(ctx context.Context, commentID, userID int64) (bool, error) {
	return false, nil
}
func (r *testCommentRepo) DeleteLikeComment(ctx context.Context, commentID, userID int64) error {
	return nil
}
func (r *testCommentRepo) StoreLikeComment(ctx context.Context, model *model.CommentLikeModel) error {
	return nil
}
func (r *testCommentRepo) GetCommentByPostIDs(ctx context.Context, postIDs []int64) ([]model.CommentModel, error) {
	return nil, nil
}
func (r *testCommentRepo) DeleteCommentsByPostID(ctx context.Context, postID int64) error {
	if r.deleteFn == nil {
		return nil
	}
	return r.deleteFn(ctx, postID)
}
