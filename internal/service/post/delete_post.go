package post

import (
	"context"
	"errors"
	"net/http"
	"time"
)

func (s *postService) DeletePost(ctx context.Context, postID, userID int64) (int, error) {
	// check post was exists
	postExists, err := s.postRepo.GetPostById(ctx, postID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if postExists == nil {
		return http.StatusNotFound, errors.New("post not found")
	}

	if postExists.UserID != userID {
		return http.StatusNotFound, errors.New("post not found")
	}

	// delete related comments and likes before soft-deleting post
	err = s.commentRepo.DeleteCommentsByPostID(ctx, postID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// soft delete post
	err = s.postRepo.SoftDelete(ctx, postID, time.Now())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// return
	return http.StatusOK, nil
}
