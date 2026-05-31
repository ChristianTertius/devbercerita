package user

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"database/sql"
	"time"
)

func (r *userRepository) GetRefreshToken(ctx context.Context, userID int64, now time.Time) (*model.RefreshTokenModel, error) {
	query := `select id, user_id, refresh_token, expired_at, created_at, updated_at from refresh_token where user_id = ? and expired_at >= ?`

	row := r.db.QueryRowContext(ctx, query, userID, now)
	var result model.RefreshTokenModel
	err := row.Scan(&result.ID, &result.UserID, &result.RefreshToken, &result.ExpiredAt, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}
