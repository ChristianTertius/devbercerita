package post

import (
	"context"
	"errors"
	"time"
)

func (r *postRepository) SoftDelete(ctx context.Context, postID int64, now time.Time) error {
	query := `update posts set deleted_at = ? where id = ?`

	result, err := r.db.ExecContext(ctx, query, now, postID)

	if err != nil {
		return err
	}

	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		return errors.New("nothing to update data")
	}

	return nil
}
