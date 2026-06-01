package post

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"database/sql"
)

func (r *postRepository) GetPostById(ctx context.Context, postID int64) (*model.PostWithUserModel, error) {
	query := `select p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at, u.username, COUNT(pl.id) as like_count
    from posts as p
    join users as u on p.user_id = u.id
    left join post_likes as pl on pl.post_id = p.id
    where p.id = ?
    and deleted_at is null
    group by p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at, u.username`

	row := r.db.QueryRowContext(ctx, query, postID)
	var result model.PostWithUserModel
	err := row.Scan(&result.ID, &result.Title, &result.Content, &result.UserID, &result.CreatedAt, &result.UpdatedAt, &result.Username, &result.LikeCount)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}
