package post

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"errors"
	"net/http"
	"time"
)

func (s *postService) LikeOrUnlikePost(ctx context.Context, postID, userID int64) (int, error) {
	// check post was exists
	postExist, err := s.postRepo.GetPostById(ctx, postID, 0)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if postExist == nil {
		return http.StatusNotFound, errors.New("post not found!")
	}

	// check user already like or not
	isUserAlreadyLike, err := s.postRepo.IsUserAlreadyLikePost(ctx, postID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// if user already like, delete like data
	if isUserAlreadyLike {
		err = s.postRepo.DeleteLikePost(ctx, postID, userID)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	} else {
		// store data
		now := time.Now()
		err := s.postRepo.StoreLikePost(ctx, &model.PostLikeModel{
			UserID:    userID,
			PostID:    postID,
			CreatedAt: now,
			UpdatedAt: now,
		})

		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	// return
	return http.StatusOK, nil
}
