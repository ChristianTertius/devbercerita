package post

import "context"

func (r *postRepository) TotalPost(ctx context.Context, search string) (int64, error) {
	query := `
        SELECT COUNT(p.id)
        FROM posts AS p
        WHERE p.deleted_at IS NULL`

	args := []any{}
	if search != "" {
		query += " AND (p.title LIKE ? OR p.content LIKE ?)"
		like := "%" + search + "%"
		args = append(args, like, like)
	}

	var total int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	return total, err
}
