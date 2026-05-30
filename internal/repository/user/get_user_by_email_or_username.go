package user

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"database/sql"
)

func (r *userRepository) GetUserByEmailOrUsername(ctx context.Context, email, username string) (*model.UserModel, error) {
	var (
		query string
		args  []interface{}
	)

	if email != "" && username != "" {
		query = `select id, username, email, password, created_at, updated_at from users where email = ? or username = ?`
		args = []interface{}{email, username}
	} else if email != "" {
		query = `select id, username, email, password, created_at, updated_at from users where email = ?`
		args = []interface{}{email}
	} else if username != "" {
		query = `select id, username, email, password, created_at, updated_at from users where username = ?`
		args = []interface{}{username}
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	var result model.UserModel
	err := row.Scan(&result.ID, &result.Username, &result.Email, &result.Password, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}
