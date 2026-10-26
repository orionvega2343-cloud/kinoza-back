package domain

import "errors"

var (
	ErrNotFound        = errors.New("watch progress not found")
	ErrInvalidDuration = errors.New("duration must be greater than zero")
	ErrInvalidPosition = errors.New("position must be between zero and duration")
	ErrDatabase        = errors.New("history error from database")
)
