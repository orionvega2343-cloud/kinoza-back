package middlewares

import (
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery - иерархически вызывается первым в списке мидлваров,
// используется для отлова паники нижних вызовов, если любой нижний метод паникует,
// go начинает разворачивать стек снизу вверх, вызывает отложенную функцию,
// recover() отлавливает панику и логирует ее, не роняя все приложение
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered",
					"panic", rec,
					"stack", string(debug.Stack()),
					"method", c.Request.Method,
					"path", c.Request.URL.Path)
				c.JSON(500, gin.H{})
			}
		}()
		c.Next()
	}
}
