package comment

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"database/sql"
)

func (r *commentRepository) DetailComment(ctx context.Context, commentID int64) (*model.CommentModel, error) {
	query := `select id, post_id, user_id, content, created_at, updated_at from comments where id = ?`
	row := r.db.QueryRowContext(ctx, query, commentID)
	var result model.CommentModel
	err := row.Scan(&result.ID, &result.PostID, &result.UserID, &result.Content, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}
