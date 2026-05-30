package main

import (
	"ChristianTertius/devbercerita/internal/config"
	"ChristianTertius/devbercerita/internal/handler/user"
	userRepo "ChristianTertius/devbercerita/internal/repository/user"
	userService "ChristianTertius/devbercerita/internal/service/user"
	"ChristianTertius/devbercerita/pkg/internalsql"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
	r := gin.Default()
	validate := validator.New()

	cfg, err := config.LoadConfig()

	db, err := internalsql.ConnectMysql(cfg)
	if err != nil {
		log.Fatal(err)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "it's work! berjalan yeyeyey",
		})
	})

	userRepo := userRepo.NewRepository(db)

	userService := userService.NewService(cfg, userRepo)

	userHandler := user.NewHandler(r, validate, userService)

	userHandler.RouteList(cfg.SecretJwt)

	addr := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("server listening on :%s", addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen error: %v", err)
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
