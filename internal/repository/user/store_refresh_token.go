package user

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
)

func (r *userRepository) StoreRefreshToken(ctx context.Context, model *model.RefreshTokenModel) error {
	query := `insert into refresh_token(user_id, refresh_token, expired_at, created_at, updated_at) values (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, model.UserID, model.RefreshToken, model.ExpiredAt, model.CreatedAt, model.UpdatedAt)
	return err
}
