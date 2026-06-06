package post

import (
	"ChristianTertius/devbercerita/internal/dto"
	"context"
	"errors"
	"net/http"
)

func (s *postService) DetailPost(ctx context.Context, postID int64, userID int64) (*dto.DetailPostResponse, int, error) {
	post, err := s.postRepo.GetPostById(ctx, postID, userID) // ← pass userID
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if post == nil {
		return nil, http.StatusNotFound, errors.New("post not found!")
	}

	postIDs := []int64{postID}
	comments, err := s.commentRepo.GetCommentByPostIDs(ctx, postIDs, userID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	commentMap := make([]dto.Comment, 0)
	for _, comment := range comments {
		commentMap = append(commentMap, dto.Comment{
			ID:        comment.ID,
			Username:  comment.Username,
			Content:   comment.Content,
			LikeCount: comment.LikeCount,
			IsLiked:   comment.IsLiked,
			CreatedAt: comment.CreatedAt.String(),
			UpdatedAt: comment.UpdatedAt.String(),
		})
	}

	return &dto.DetailPostResponse{
		ID:        post.ID,
		Username:  post.Username,
		Title:     post.Title,
		Content:   post.Content,
		LikeCount: post.LikeCount,
		IsLiked:   post.IsLiked, // ← tambah
		Comments:  commentMap,
		CreatedAt: post.CreatedAt.String(),
		UpdatedAt: post.UpdatedAt.String(),
	}, http.StatusOK, nil
}
