# DevBercerita Backend

A lightweight RESTful backend written in Go, powered by Gin, that handles account lifecycle, posts, comments, and social interactions (likes) backed by MySQL and a JWT authentication layer.

## Features
- User registration, login, and refresh-token rotation guarded by a refresh-token middleware.
- CRUD for posts with soft deletes plus like/unlike toggles.
- Commenting system with its own like/unlike actions and post-scoped retrieval.
- Layered architecture (handler → service → repository) with DTOs and models to keep behavior testable and boundaries explicit.

## Architecture overview
- **Gin** exposes HTTP endpoints defined in `internal/handler/*` and reuses shared validation/middleware helpers.
- **Services** encapsulate business rules (`internal/service/*`). They orchestrate repositories, JWT helpers, and DTO mapping.
- **Repositories** talk to MySQL via `pkg/internalsql` and live under `internal/repository/*` (posts, comments, users, likes, refresh tokens).
- **Models/DTOs** keep payloads and persistence consistent (`internal/model/*`, `internal/dto/*`).
- **Middleware** enforces access control plus refresh token validation (`internal/middleware/middleware.go`).

## Getting started

### Prerequisites
- Go 1.26.3 (matches `go.mod`).
- MySQL 8.x (Docker compose is provided).
- `dbmate` or raw SQL client if you plan to run migrations manually.

### Environment configuration
Copy `.env.example` to `.env` (or export the values directly). Key values the service reads:

| Variable | Purpose |
| --- | --- |
| `PORT` | HTTP port |
| `DATABASE_URL` | Full DSN for migrations (optional) |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | Used by the connection helper |
| `SECRET_JWT` | HMAC signing secret for access tokens (required) |
| `app_env` | `development`/`production` toggles logging |

### Database setup
1. Spin up MySQL via Docker: `docker compose up -d db` (binds host port `3307`).
2. Apply migrations located in `db/migrations/` (either with `dbmate` or by executing `db/schema.sql`).
3. Ensure the `refresh_token`, `post_likes`, `comment_likes`, and related tables exist before running the server.

### Run the service

```bash
go run ./cmd/main.go
```

The server listens on `:PORT` (default `8080`) and exposes a `/health` endpoint for quick availability checks.

## API endpoints

### Authentication
- `POST /auth/register`: register with `email`, `username`, `password`.
- `POST /auth/login`: obtain access + refresh tokens.
- `POST /auth/refresh`: rotate refresh token (requires refresh token middleware).

### Posts (`/posts`)
- `POST /posts`: create (auth required).
- `PUT /posts/:post_id/update`: update owned post.
- `DELETE /posts/:post_id/delete`: soft delete.
- `POST /posts/action`: like or unlike a post.
- `GET /posts/:post_id/detail`: fetch post with comments/likes metadata.
- `GET /posts/`: paginated feed.

### Comments (`/comments`)
- `POST /comments`: attach a comment to a post.
- `POST /comments/action`: like/unlike a comment.

Authentication is JWT-based; supply `Authorization: Bearer <token>` for protected routes.

## Project layout
- `cmd/main.go` – entry point, wiring config, repositories, services, and handlers.
- `internal/config` – loads `.env` (with gin defaults) and centralizes the config struct.
- `pkg/jwt` & `pkg/refreshtoken` – helpers to mint/validate tokens.
- `pkg/internalsql` – wrapper around the MySQL driver for consistent connections.
