package post

import "context"

func (r *postRepository) TotalPost(ctx context.Context) (int64, error) {
	query := `select count(id) from posts where deleted_at is null`

	var total int64
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
