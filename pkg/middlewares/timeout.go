package middlewares

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout - используется, для предотвращения зависания операции, к примеру SQL запрос не отвечает,
// timeout идет вверх по стеку, и передает информацию в transaction,
// которая в свою очередь откатывает метод
func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
