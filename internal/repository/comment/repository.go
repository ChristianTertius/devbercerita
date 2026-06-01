package comment

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"database/sql"
)

type CommentRepository interface {
	StoreComment(ctx context.Context, model *model.CommentModel) error
}

type commentRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}
