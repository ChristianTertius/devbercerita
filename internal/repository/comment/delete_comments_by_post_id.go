package comment

import (
	"context"
)

func (r *commentRepository) DeleteCommentsByPostID(ctx context.Context, postID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if _, err = tx.ExecContext(ctx, `DELETE cl FROM comment_likes cl JOIN comments c ON cl.comment_id = c.id WHERE c.post_id = ?`, postID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM comments WHERE post_id = ?`, postID); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
