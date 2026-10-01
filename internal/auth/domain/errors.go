package domain

import "errors"

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrInvalidEmail = errors.New("email is empty")
var ErrInvalidRole = errors.New("role is empty")
var ErrTokenExpired = errors.New("token expired")
var ErrTokenRevoked = errors.New("token already revoked")
