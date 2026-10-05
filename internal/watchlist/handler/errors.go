package handler

import (
	"errors"

	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/pkg/response"
)

func mapServiceError(err error) (int, response.Error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return 404, response.Error{Message: err.Error(), Code: "NOT_FOUND"}
	case errors.Is(err, domain.ErrAlreadyExists):
		return 409, response.Error{Message: err.Error(), Code: "ALREADY_EXISTS"}
	case errors.Is(err, domain.ErrInvalidStatus):
		return 400, response.Error{Message: err.Error(), Code: "INVALID_STATUS"}
	case errors.Is(err, domain.ErrInvalidStatus):
		return 400, response.Error{Message: err.Error(), Code: "INVALID_STATUS"}
	default:
		return 500, response.Error{Message: "internal server error", Code: "INTERNAL_ERROR"}
	}
}
