package post

import (
	"ChristianTertius/devbercerita/internal/dto"
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"errors"
	"net/http"
	"time"
)

func (s *postService) UpdatePost(ctx context.Context, req *dto.CreateOrUpdatePostRequest, postID, userID int64) (int, error) {
	// check if post was exist
	postExists, err := s.postRepo.GetPostById(ctx, postID, 0)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if postExists == nil {
		return http.StatusNotFound, errors.New("post not found!")
	}

	if postExists.UserID != userID {
		return http.StatusNotFound, errors.New("post not found!")
	}

	// update the post
	err = s.postRepo.UpdatePost(ctx, &model.PostModel{
		Title:     req.Title,
		Content:   req.Content,
		UpdatedAt: time.Now(),
	}, postID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// return
	return http.StatusOK, nil
}
