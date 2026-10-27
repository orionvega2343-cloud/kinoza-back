package domain

import "errors"

var (
	ErrInvalidPosition = errors.New("position must not be negative")
	ErrDatabase        = errors.New("history error from database")
)
