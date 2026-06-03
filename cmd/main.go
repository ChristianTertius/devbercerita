// Package main bootstraps the DevBercerita HTTP API.
// @title DevBercerita API
// @version 1.0
// @description API untuk platform DevBercerita: register/login, posting, dan komentar.
// @contact.name DevBercerita Team
// @contact.email hi@devbercerita.local
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"ChristianTertius/devbercerita/docs"
	"ChristianTertius/devbercerita/internal/config"
	commentHandler "ChristianTertius/devbercerita/internal/handler/comment"
	postHandler "ChristianTertius/devbercerita/internal/handler/post"
	"ChristianTertius/devbercerita/internal/handler/user"
	commentRepo "ChristianTertius/devbercerita/internal/repository/comment"
	postRepo "ChristianTertius/devbercerita/internal/repository/post"
	userRepo "ChristianTertius/devbercerita/internal/repository/user"
	commentService "ChristianTertius/devbercerita/internal/service/comment"
	postService "ChristianTertius/devbercerita/internal/service/post"
	userService "ChristianTertius/devbercerita/internal/service/user"
	"ChristianTertius/devbercerita/pkg/internalsql"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()
	validate := validator.New()

	cfg, err := config.LoadConfig()

	db, err := internalsql.ConnectMysql(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// config origin
	origins := os.Getenv("ALLOWED_ORIGINS")
	allowOrigins := strings.Split(origins, ",")

	// local dev fallback
	if len(allowOrigins) == 0 || allowOrigins[0] == "" {
		allowOrigins = []string{"http://127.0.0.1:3000", "http://localhost:3000"}
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "it's work",
		})
	})

	userRepo := userRepo.NewRepository(db)
	postRepo := postRepo.NewPostRepository(db)
	commentRepo := commentRepo.NewRepository(db)

	userService := userService.NewService(cfg, userRepo)
	postService := postService.NewService(cfg, postRepo, commentRepo)
	commentService := commentService.NewCommentService(cfg, commentRepo, postRepo)

	userHandler := user.NewHandler(r, validate, userService)
	postHandler := postHandler.NewHandler(r, validate, postService)
	commentHandler := commentHandler.NewHandler(r, validate, commentService)

	userHandler.RouteList(cfg.SecretJwt)
	postHandler.RouteList(cfg.SecretJwt)
	commentHandler.RouteList(cfg.SecretJwt)

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Port)
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("server listening on :%s", addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// tunggu sinyal dari os
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("Shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("forced shutdown: %v", err)
	}

	log.Println("server stopped")
}
