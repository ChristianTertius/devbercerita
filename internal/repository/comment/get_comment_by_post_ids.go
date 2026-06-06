package comment

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (r *commentRepository) GetCommentByPostIDs(ctx context.Context, postIDs []int64, userID int64) ([]model.CommentModel, error) {
	if len(postIDs) == 0 {
		return []model.CommentModel{}, nil
	}
	placeholders := make([]string, len(postIDs))
	args := make([]interface{}, len(postIDs)+1)
	args[0] = userID
	for i, id := range postIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}
	query := fmt.Sprintf(`SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at, c.updated_at,
        COUNT(cl.id) as like_count,
        EXISTS(SELECT 1 FROM comment_likes WHERE comment_id = c.id AND user_id = ?) as is_liked
        FROM comments as c
        JOIN users as u ON u.id = c.user_id
        LEFT JOIN comment_likes as cl ON cl.comment_id = c.id
        WHERE c.post_id IN (%s)
        GROUP BY c.id, c.post_id, c.user_id, u.username, c.content, c.created_at, c.updated_at
        ORDER BY like_count DESC`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.CommentModel{}, nil
		}
		return []model.CommentModel{}, err
	}
	result := make([]model.CommentModel, 0)
	for rows.Next() {
		var data model.CommentModel
		err = rows.Scan(&data.ID, &data.PostID, &data.UserID, &data.Username, &data.Content,
			&data.CreatedAt, &data.UpdatedAt, &data.LikeCount, &data.IsLiked)
		if err != nil {
			return []model.CommentModel{}, err
		}
		result = append(result, data)
	}
	return result, nil
}
