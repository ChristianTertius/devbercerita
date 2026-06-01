package post

import "context"

func (r *postRepository) DeleteLikePost(ctx context.Context, postID, userID int64) error {
	query := `delete from post_likes where post_id = ? and user_id = ?`
	_, err := r.db.ExecContext(ctx, query, postID, userID)
	return err
}
