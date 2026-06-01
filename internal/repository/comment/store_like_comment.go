package comment

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
)

func (r *commentRepository) StoreLikeComment(ctx context.Context, model *model.CommentLikeModel) error {
	query := `insert into comment_likes (comment_id, user_id, created_at, updated_at) values (?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, model.CommentID, model.UserID, model.CreatedAt, model.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}
