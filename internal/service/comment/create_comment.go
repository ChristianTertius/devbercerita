package comment

import (
	"ChristianTertius/devbercerita/internal/dto"
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"errors"
	"net/http"
	"time"
)

func (s *commentService) CreateComment(ctx context.Context, req *dto.StoreCommentRequest, userID int64) (int, error) {
	// check if post was exists
	postExist, err := s.postRepo.GetPostById(ctx, req.PostID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if postExist == nil {
		return http.StatusNotFound, errors.New("Post not found!")
	}

	// store comment
	now := time.Now()
	err = s.commentRepo.StoreComment(ctx, &model.CommentModel{
		PostID:    req.PostID,
		UserID:    userID,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	})

	if err != nil {
		return http.StatusInternalServerError, err
	}

	// return
	return http.StatusOK, nil
}
