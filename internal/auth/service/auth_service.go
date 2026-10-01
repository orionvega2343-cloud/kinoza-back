package service

import (
	cl "kinoza-back/pkg/claims"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func newAccessToken(userId, email, role, secret string, ttl time.Duration) (string, error) {
	claims := cl.Claims{
		UserId: userId,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
