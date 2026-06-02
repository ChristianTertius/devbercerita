package comment

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	_ "github.com/go-sql-driver/mysql"
)

const schemaStatements = `
DROP TABLE IF EXISTS comment_likes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
	id INT NOT NULL AUTO_INCREMENT,
	email VARCHAR(255) NOT NULL,
	username VARCHAR(255) NOT NULL,
	password VARCHAR(500) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE posts (
	id INT NOT NULL AUTO_INCREMENT,
	user_id INT NOT NULL,
	title VARCHAR(255) NOT NULL,
	content LONGTEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP NULL DEFAULT NULL,
	PRIMARY KEY (id),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE comments (
	id INT NOT NULL AUTO_INCREMENT,
	post_id INT NOT NULL,
	user_id INT NOT NULL,
	content LONGTEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE comment_likes (
	id INT NOT NULL AUTO_INCREMENT,
	comment_id INT NOT NULL,
	user_id INT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`

func TestDeleteCommentsByPostIDIntegration(t *testing.T) {
	withMySQL(t, func(db *sql.DB) {
		repo := &commentRepository{db: db}
		if err := seedPostWithComments(db); err != nil {
			t.Fatalf("seed failed: %v", err)
		}

		if err := repo.DeleteCommentsByPostID(context.Background(), 1); err != nil {
			t.Fatalf("delete failed: %v", err)
		}

		assertCount(t, db, "comments", 0)
		assertCount(t, db, "comment_likes", 0)
	})
}

func withMySQL(t *testing.T, fn func(db *sql.DB)) {
	t.Helper()
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Skipf("docker not available: %v", err)
		return
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.1",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=secret",
			"MYSQL_DATABASE=testdb",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"3306/tcp": {{HostIP: "0.0.0.0", HostPort: ""}},
		},
	})
	if err != nil {
		t.Fatalf("container start failed: %v", err)
	}
	resource.Expire(900)
	pool.MaxWait = 2 * time.Minute
	defer func() {
		_ = pool.Purge(resource)
	}()

	port := resource.GetPort("3306/tcp")
	dsn := fmt.Sprintf("root:secret@(localhost:%s)/testdb?parseTime=true&multiStatements=true", port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open db failed: %v", err)
	}
	defer db.Close()

	if err := pool.Retry(func() error {
		return db.Ping()
	}); err != nil {
		t.Fatalf("ping failed: %v", err)
	}

	for _, stmt := range strings.Split(schemaStatements, ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("schema exec failed: %v", err)
		}
	}

	fn(db)
}

func seedPostWithComments(db *sql.DB) error {
	if _, err := db.Exec(`INSERT INTO users (email, username, password) VALUES (?, ?, ?)`, "a@a.com", "user", "pass"); err != nil {
		return err
	}

	if _, err := db.Exec(`INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)`, 1, "title", "body"); err != nil {
		return err
	}

	if _, err := db.Exec(`INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`, 1, 1, "comment"); err != nil {
		return err
	}

	if _, err := db.Exec(`INSERT INTO comment_likes (comment_id, user_id) VALUES (?, ?)`, 1, 1); err != nil {
		return err
	}

	return nil
}

func assertCount(t *testing.T, db *sql.DB, table string, want int) {
	t.Helper()
	var got int
	if err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&got); err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	if got != want {
		t.Fatalf("expected %d rows in %s, got %d", want, table, got)
	}
}
