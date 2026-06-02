package post

import (
	"ChristianTertius/devbercerita/internal/dto"
	"ChristianTertius/devbercerita/internal/model"
	"context"
)

func (r *postRepository) GetAllPost(ctx context.Context, param *dto.GetAllPostRequest, offset int) ([]model.PostWithUserModel, error) {
	query := `select 
    p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at, u.username, count(pl.id) as like_count
    from posts as p
    join users as u on u.id = p.user_id
    left join post_likes as pl on pl.post_id = p.id
    where p.deleted_at is null
    group by p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at, u.username
    order by created_at desc
    limit ?
    offset ? `

	rows, err := r.db.QueryContext(ctx, query, param.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.PostWithUserModel, 0)

	for rows.Next() {
		var data model.PostWithUserModel
		err := rows.Scan(&data.ID, &data.Title, &data.Content, &data.UserID, &data.CreatedAt, &data.UpdatedAt, &data.Username, &data.LikeCount)
		if err != nil {
			return nil, err
		}

		result = append(result, model.PostWithUserModel{
			ID:        data.ID,
			UserID:    data.UserID,
			Title:     data.Title,
			Content:   data.Content,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			Username:  data.Username,
			LikeCount: data.LikeCount,
		})
	}

	return result, nil
}
