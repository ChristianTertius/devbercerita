package comment

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// or GetAllComments
func (r *commentRepository) GetCommentByPostIDs(ctx context.Context, postIDs []int64) ([]model.CommentModel, error) {
	if len(postIDs) == 0 {
		return []model.CommentModel{}, nil
	}
	placeholders := make([]string, len(postIDs))
	args := make([]interface{}, len(postIDs))
	for i, id := range postIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`select c.id, c.post_id, c.user_id, u.username, c.content, c.created_at, c.updated_at, count(cl.id) as like_count 
        from comments as c
        join users as u on u.id = c.user_id
        left join comment_likes as cl on cl.comment_id = c.id
        where c.post_id in (%s)
        group by c.id, c.post_id, c.user_id, u.username, c.content, c.created_at, c.updated_at
        order by like_count desc`, strings.Join(placeholders, ","))

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
		err = rows.Scan(&data.ID, &data.PostID, &data.UserID, &data.Username, &data.Content, &data.CreatedAt, &data.UpdatedAt, &data.LikeCount)
		if err != nil {
			return []model.CommentModel{}, err
		}

		result = append(result, model.CommentModel{
			ID:        data.ID,
			PostID:    data.PostID,
			UserID:    data.UserID,
			Username:  data.Username,
			Content:   data.Content,
			LikeCount: data.LikeCount,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
		})
	}
	return result, nil
}
