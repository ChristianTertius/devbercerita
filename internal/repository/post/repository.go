package post

import (
	"ChristianTertius/devbercerita/internal/dto"
	"ChristianTertius/devbercerita/internal/model"
	"context"
	"database/sql"
	"time"
)

type PostRepository interface {
	StorePost(ctx context.Context, model *model.PostModel) (int64, error)
	GetPostById(ctx context.Context, postID int64, userID int64) (*model.PostWithUserModel, error)
	UpdatePost(ctx context.Context, model *model.PostModel, postID int64) error
	SoftDelete(ctx context.Context, postID int64, now time.Time) error
	IsUserAlreadyLikePost(ctx context.Context, postID, userID int64) (bool, error)
	DeleteLikePost(ctx context.Context, postID, userID int64) error
	StoreLikePost(ctx context.Context, model *model.PostLikeModel) error
	TotalPost(ctx context.Context, search string) (int64, error)
	GetAllPost(ctx context.Context, param *dto.GetAllPostRequest, offset int) ([]model.PostWithUserModel, error)
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}
