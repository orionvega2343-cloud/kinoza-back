package middlewares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger - пишет структурированный лог по каждому запросу: method,
// path, query, start, status, clientIp
// для логирования поведения приложения, статус по умолчанию 200,
// даже если handler не вызвал WriteHeader, go все равно вернет 200,
// записывает в журнал все действия нижних обработчиков
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if query != "" {
			path = path + "?" + query
		}

		ctx := IDFromContext(c.Request.Context())
		slog.Info("http",
			"status", statusCode,
			"latency", latency,
			"clientIP", clientIP,
			"method", method,
			"path", path,
			"query", query,
			"requestId", ctx)
	}
}
