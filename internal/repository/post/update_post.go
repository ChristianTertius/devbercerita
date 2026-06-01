package post

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"errors"
)

func (r *postRepository) UpdatePost(ctx context.Context, model *model.PostModel, postID int64) error {
	query := `update posts set title = ?, content = ?, updated_at = ? where id = ?`
	result, err := r.db.ExecContext(ctx, query, model.Title, model.Content, model.UpdatedAt, postID)
	if err != nil {
		return err
	}

	rowAffected, _ := result.RowsAffected()
	if rowAffected == 0 {
		return errors.New("nothing to update")
	}

	return nil
}
