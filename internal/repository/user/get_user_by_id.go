package user

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"database/sql"
)

func (r *userRepository) GetUserByID(ctx context.Context, userID int64) (*model.UserModel, error) {
	query := `select id, username, email, created_at, updated_at from users where id = ?`

	row := r.db.QueryRowContext(ctx, query, userID)
	var result model.UserModel

	err := row.Scan(&result.ID, &result.Username, &result.Email, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}
