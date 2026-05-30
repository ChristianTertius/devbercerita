package user

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
)

func (r *userRepository) CreateUser(ctx context.Context, model *model.UserModel) (int64, error) {
	query := `insert into users(email, username, password, created_at, updated_at) values (?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, model.Email, model.Username, model.Password, model.CreatedAt, model.UpdatedAt)
	if err != nil {
		return 0, err
	}

	userID, _ := result.LastInsertId()

	return userID, nil
}
