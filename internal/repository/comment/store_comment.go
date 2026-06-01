package comment

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
)

func (r *commentRepository) StoreComment(ctx context.Context, model *model.CommentModel) error {
	query := `insert into comments (post_id, user_id, content, created_at, updated_at) values (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, model.PostID, model.UserID, model.Content, model.CreatedAt, model.UpdatedAt)

	return err
}
