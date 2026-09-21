package main

import (
	"kinoza-back/pkg/config"
	"kinoza-back/pkg/logger"
	"kinoza-back/pkg/postgres"
	"log/slog"
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

}
