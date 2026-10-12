package domain

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("watchlist item not found")
	ErrInvalidStatus = errors.New("watchlist item has an invalid status")
)
