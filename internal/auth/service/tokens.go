package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	cl "kinoza-back/pkg/claims"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// newAccessToken - внутри себя собирает Claims, создает новый токен
// с подписью HS256, формирует и возвращает подписанный JWT Токен
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

// hashToken - хэширует переданную строку (sha256 + hex),
// используется как для новых токенов, так и для сверки токена,
// присланного клиентом
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// newRefreshToken - создает последовательный набор байт,
// заполняет слайс key криптографически стойкими рандомными байтами,
// идиоматично не возвращает ошибку, кроме как в старых Linux системах,
// для оптимизации было принято решение ее обработать, хэширует контрольную сумму,
// переводит key и sum в строки и возвращает их
func newRefreshToken() (raw string, hash string, err error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		slog.Error("failed to generate random key", "error", err)
		return "", "", err
	}
	raw = hex.EncodeToString(key)
	return raw, hashToken(raw), nil
}
