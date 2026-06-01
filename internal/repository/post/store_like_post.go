package post

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
)

func (r *postRepository) StoreLikePost(ctx context.Context, model *model.PostLikeModel) error {
	query := `insert into post_likes (post_id, user_id, created_at, updated_at) values (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, model.PostID, model.UserID, model.CreatedAt, model.UpdatedAt)
	return err
}
