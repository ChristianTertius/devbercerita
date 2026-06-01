package comment

import (
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"database/sql"
)

type CommentRepository interface {
	StoreComment(ctx context.Context, model *model.CommentModel) error
	DetailComment(ctx context.Context, commentID int64) (*model.CommentModel, error)
	IsUserAlreadyLikeComment(ctx context.Context, commentID, userID int64) (bool, error)
	DeleteLikeComment(ctx context.Context, commentID, userID int64) error
	StoreLikeComment(ctx context.Context, model *model.CommentLikeModel) error
}

type commentRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}
