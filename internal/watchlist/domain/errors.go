package domain

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("watchlist item not found")
	ErrAlreadyExists = errors.New("watchlist item already exists")
)
