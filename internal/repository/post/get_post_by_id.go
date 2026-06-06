package post

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"database/sql"
)

func (r *postRepository) GetPostById(ctx context.Context, postID int64, userID int64) (*model.PostWithUserModel, error) {
	query := `SELECT p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at, u.username,
    COUNT(pl.id) as like_count,
    EXISTS(SELECT 1 FROM post_likes WHERE post_id = p.id AND user_id = ?) as is_liked
    FROM posts as p
    JOIN users as u ON p.user_id = u.id
    LEFT JOIN post_likes as pl ON pl.post_id = p.id
    WHERE p.id = ?
    AND deleted_at IS NULL
    GROUP BY p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at, u.username`

	row := r.db.QueryRowContext(ctx, query, userID, postID)
	var result model.PostWithUserModel
	err := row.Scan(&result.ID, &result.Title, &result.Content, &result.UserID,
		&result.CreatedAt, &result.UpdatedAt, &result.Username, &result.LikeCount, &result.IsLiked)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
