package internalsql

import (
	"ChristianTertius/devbercerita/internal/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectMysql(cfg *config.Config) (*sql.DB, error) {
	dsn, err := buildDSN(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.AppEnv == "production" {
		dsn += "&tls=skip-verify"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("database connected")
	return db, nil
}

func buildDSN(cfg *config.Config) (string, error) {
	if cfg.DBUrlMigration != "" {
		raw := cfg.DBUrlMigration
		if !strings.Contains(raw, "?") {
			raw += ""
		}

		u, err := url.Parse(raw)
		if err != nil {
			return "", fmt.Errorf("invalid DATABASE_URL: %w", err)
		}

		user := u.User.Username()
		password, _ := u.User.Password()
		host := u.Hostname()
		port := u.Port()
		if port == "" {
			port = "3306"
		}
		db := strings.Trim(u.Path, "/")
		if db == "" {
			return "", fmt.Errorf("invalid DATABASE_URL: missing database name")
		}

		params := url.Values{}
		params.Set("parseTime", "true")
		params.Set("timeout", "5s")
		params.Set("readTimeout", "5s")
		if u.RawQuery != "" {
			for key, values := range u.Query() {
				if len(values) > 0 {
					params.Set(key, values[0])
				}
			}
		}

		if user == "" {
			return "", fmt.Errorf("invalid DATABASE_URL: missing username")
		}

		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, password, host, port, db, params.Encode()), nil
	}

	params := url.Values{}
	params.Set("parseTime", "true")
	params.Set("timeout", "5s")
	params.Set("readTimeout", "5s")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, params.Encode()), nil
}
