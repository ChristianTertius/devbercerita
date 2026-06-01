package comment

import "context"

func (r *commentRepository) DeleteLikeComment(ctx context.Context, commentID, userID int64) error {
	query := `delete from comment_likes where comment_id = ? and user_id = ?`
	_, err := r.db.ExecContext(ctx, query, commentID, userID)
	if err != nil {
		return err
	}

	return nil
}
