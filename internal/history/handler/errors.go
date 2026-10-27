package handler

import (
	"errors"

	"kinoza-back/internal/history/domain"
	"kinoza-back/pkg/response"
)

func mapServiceError(err error) (int, response.Error) {
	switch {
	case errors.Is(err, domain.ErrInvalidPosition):
		return 400, response.Error{Message: err.Error(), Code: "INVALID_POSITION"}
	default:
		return 500, response.Error{Message: "internal server error", Code: "INTERNAL_ERROR"}
	}
}
