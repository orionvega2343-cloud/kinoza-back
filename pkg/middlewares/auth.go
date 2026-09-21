package middlewares

import (
	"fmt"
	cl "kinoza-back/pkg/claims"
	"kinoza-back/pkg/response"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth - middleware для проверки JWT токена,
// ожидает заголовок Authorization: Bearer <token>,
// при успешной валидации кладёт user_id в контекст запроса,
// при ошибке — прерывает цепочку и возвращает 401
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, response.Error{Message: "no token provided", Code: "NO_TOKEN_PROVIDED"})
			return
		}
		//Проверяем формат Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(401, response.Error{Message: "invalid token", Code: "INVALID_TOKEN"})
			return
		}
		tokenString := parts[1]
		claims := cl.Claims{}
		//парсинг и валидация токена
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_KEY")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}
		c.Set("role", claims.Role)
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
