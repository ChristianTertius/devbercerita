package post

import (
	"ChristianTertius/devbercerita/internal/dto"
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"fmt"
)

func (r *postRepository) GetAllPost(ctx context.Context, param *dto.GetAllPostRequest, offset int) ([]model.PostWithUserModel, error) {
	query := fmt.Sprintf(`
        SELECT
            p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at,
            u.username, COUNT(pl.id) AS like_count
        FROM posts AS p
        JOIN users AS u ON u.id = p.user_id
        LEFT JOIN post_likes AS pl ON pl.post_id = p.id
        WHERE p.deleted_at IS NULL
        %s
        GROUP BY p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at, u.username
        ORDER BY %s %s
        LIMIT ? OFFSET ?`,
		searchClause(param.Search),
		sanitizeSortBy(param.SortBy),
		sanitizeOrder(param.Order),
	)

	args := buildArgs(param.Search, param.Limit, int64(offset))
	rows, err := r.db.QueryContext(ctx, query, args...)
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

func searchClause(search string) string {
	if search == "" {
		return ""
	}
	return "AND (p.title LIKE ? OR p.content LIKE ?)"
}

func buildArgs(search string, limit, offset int64) []any {
	if search != "" {
		like := "%" + search + "%"
		return []any{like, like, limit, offset}
	}
	return []any{limit, offset}
}

func sanitizeSortBy(s string) string {
	if s == "like_count" {
		return "like_count"
	}
	return "p.created_at"
}

func sanitizeOrder(o string) string {
	if o == "asc" {
		return "ASC"
	}
	return "DESC"
}
