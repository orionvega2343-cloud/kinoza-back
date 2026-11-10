package handler

import (
	"errors"
	"net/http"

	"kinoza-back/internal/friends/domain"
	"kinoza-back/pkg/response"
)

func mapServiceError(err error) (int, response.Error) {
	switch {
	case errors.Is(err, domain.ErrSelfRequest):
		return http.StatusBadRequest, response.Error{Message: err.Error(), Code: "SELF_REQUEST"}
	case errors.Is(err, domain.ErrNotRecipient):
		return http.StatusForbidden, response.Error{Message: err.Error(), Code: "NOT_RECIPIENT"}
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, response.Error{Message: err.Error(), Code: "USER_NOT_FOUND"}
	case errors.Is(err, domain.ErrRequestNotFound):
		return http.StatusNotFound, response.Error{Message: err.Error(), Code: "REQUEST_NOT_FOUND"}
	case errors.Is(err, domain.ErrFriendshipNotFound):
		return http.StatusNotFound, response.Error{Message: err.Error(), Code: "FRIENDSHIP_NOT_FOUND"}
	case errors.Is(err, domain.ErrRequestAlreadyExists):
		return http.StatusConflict, response.Error{Message: err.Error(), Code: "REQUEST_ALREADY_EXISTS"}
	case errors.Is(err, domain.ErrAlreadyFriends):
		return http.StatusConflict, response.Error{Message: err.Error(), Code: "ALREADY_FRIENDS"}
	case errors.Is(err, domain.ErrInvalidTransition):
		return http.StatusConflict, response.Error{Message: err.Error(), Code: "INVALID_TRANSITION"}
	default:
		return http.StatusInternalServerError, response.Error{Message: "internal server error", Code: "INTERNAL_ERROR"}
	}
}
