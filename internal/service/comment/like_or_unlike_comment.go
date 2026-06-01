package comment

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"errors"
	"net/http"
	"time"
)

func (s *commentService) LikeOrUnlikeComment(ctx context.Context, commentID, userID int64) (int, error) {
	// check comment is exist
	commentExist, err := s.commentRepo.DetailComment(ctx, commentID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if commentExist == nil {
		return http.StatusNotFound, errors.New("comment not found!")
	}

	// check user already like comment
	isUserAlreadyLikeComment, err := s.commentRepo.IsUserAlreadyLikeComment(ctx, commentID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if isUserAlreadyLikeComment {
		err := s.commentRepo.DeleteLikeComment(ctx, commentID, userID)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	} else {
		// store data
		now := time.Now()
		err := s.commentRepo.StoreLikeComment(ctx, &model.CommentLikeModel{
			UserID:    userID,
			CommentID: commentID,
			CreatedAt: now,
			UpdatedAt: now,
		})

		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}
