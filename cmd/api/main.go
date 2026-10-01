package main

import (
	"fmt"
	"log/slog"
	"time"

	"kinoza-back/internal/auth/handler"
	"kinoza-back/internal/auth/repository"
	"kinoza-back/internal/auth/service"
	"kinoza-back/pkg/config"
	"kinoza-back/pkg/logger"
	"kinoza-back/pkg/middlewares"
	"kinoza-back/pkg/postgres"
	"kinoza-back/pkg/transaction"

	"github.com/gin-gonic/gin"
)

func main() {
	slog.SetDefault(logger.NewLogger())
	cfg := config.MustLoad()

	db, err := postgres.Connection(*cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database connection", "error", err)
		}
	}()

	tx := transaction.NewTransactor(db)

	userRepo := repository.NewUserRepoImpl(db)
	refreshTokenRepo := repository.NewRefreshTokenRepo(db)

	cookieTTL := int(cfg.Jwt.RefreshTTL.Seconds())

	userService := service.NewUserService(
		userRepo,
		refreshTokenRepo,
		cfg.Jwt.Secret,
		cfg.Jwt.AccessTTL,
		cfg.Jwt.RefreshTTL,
		tx,
	)

	userHandler := handler.NewUserHandler(userService, cookieTTL)

	r := gin.New()
	r.Use(
		middlewares.Recovery(),
		middlewares.RequestId(),
		middlewares.Logger(),
		middlewares.Timeout(10*time.Second),
	)

	auth := r.Group("/auth")
	auth.POST("/register", userHandler.Register)
	auth.POST("/login", userHandler.Login)
	auth.POST("/refresh", userHandler.Refresh)

	protected := auth.Group("")
	protected.Use(middlewares.Auth())
	protected.GET("/users/:id", userHandler.GetUserById)
	protected.PATCH("/users/me", userHandler.UpdateUser)
	protected.GET("/users/:userId/tokens", userHandler.GetTokenList)
	protected.DELETE("/tokens/:id", userHandler.RevokeToken)
	protected.DELETE("/users/:userId/tokens", userHandler.RevokeAllTokens)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		slog.Error("server stopped", "error", err)
	}
}
