package postgres

import (
	"context"
	"fmt"
	"kinoza-back/pkg/config"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// HealthCheck - проверяет валидное соединение с БД,
// если соединения нет, установит его через PingContext
func HealthCheck(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			slog.Error("health check failed", "error", err)
			c.JSON(503, gin.H{
				"status": "failed",
				"error":  "postgres connection failed",
			})
			return
		}
		c.JSON(200, gin.H{
			"status": "ok",
		})
	}
}

func Connection(cfg config.Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", cfg.Db.Host, cfg.Db.Port, cfg.Db.User, cfg.Db.Name, cfg.Db.Password, cfg.Db.SslMode)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		slog.Error("failed connect to postgres", err)
		return nil, err
	}
	return db, nil
}
