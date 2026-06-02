package comment

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDeleteCommentsByPostID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()
	repo := &commentRepository{db: db}
	postID := int64(10)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE cl FROM comment_likes cl JOIN comments c ON cl.comment_id = c.id WHERE c.post_id = ?`).
		WithArgs(postID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`DELETE FROM comments WHERE post_id = ?`).
		WithArgs(postID).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit()

	if err := repo.DeleteCommentsByPostID(context.Background(), postID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}

func TestDeleteCommentsByPostID_RollbackOnFailure(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()
	repo := &commentRepository{db: db}
	postID := int64(20)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE cl FROM comment_likes cl JOIN comments c ON cl.comment_id = c.id WHERE c.post_id = ?`).
		WithArgs(postID).
		WillReturnError(errors.New("boom"))
	mock.ExpectRollback()

	if err := repo.DeleteCommentsByPostID(context.Background(), postID); err == nil {
		t.Fatal("expected error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations not met: %v", err)
	}
}
