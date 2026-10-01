package handler

import (
	"errors"

	"kinoza-back/internal/auth/domain"
	"kinoza-back/pkg/response"
)

// mapServiceError - переводит доменные sentinel-ошибки сервиса в HTTP-статус
// и тело ответа; для нераспознанной ошибки наружу уходит общее сообщение,
// а не err.Error() - чтобы не утекали детали внутренней реализации (например,
// текст ошибки БД)
func mapServiceError(err error) (int, response.Error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		return 401, response.Error{Message: err.Error(), Code: "INVALID_CREDENTIALS"}
	case errors.Is(err, domain.ErrTokenExpired):
		return 401, response.Error{Message: err.Error(), Code: "TOKEN_EXPIRED"}
	case errors.Is(err, domain.ErrTokenRevoked):
		return 401, response.Error{Message: err.Error(), Code: "TOKEN_REVOKED"}
	case errors.Is(err, domain.ErrInvalidEmail):
		return 400, response.Error{Message: err.Error(), Code: "INVALID_EMAIL"}
	case errors.Is(err, domain.ErrInvalidRole):
		return 400, response.Error{Message: err.Error(), Code: "INVALID_ROLE"}
	default:
		return 500, response.Error{Message: "internal server error", Code: "INTERNAL_ERROR"}
	}
}
