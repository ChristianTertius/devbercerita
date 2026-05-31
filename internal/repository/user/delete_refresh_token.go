package user

import (
	"context"
	"errors"
)

func (r *userRepository) DeleteRefreshTokenByID(ctx context.Context, userID int64) error {
	query := `delete from refresh_token where user_id = ?`

	result, err := r.db.ExecContext(ctx, query, userID)

	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return errors.New("nothing to delete! id not found!")
	}

	return nil
}
