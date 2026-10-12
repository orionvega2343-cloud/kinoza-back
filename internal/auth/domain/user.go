package domain

import (
	"log/slog"
	"strings"
	"time"
)

type User struct {
	Id           string    `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
}

func (u *User) ValidateEmail() string {
	if u.Email == "" || !strings.Contains(u.Email, "@") {
		slog.Error("email is required or the email incorrect ", "error", u.Email)
		return ""
	}
	return u.Email
}

func (u *User) ValidateRole() string {
	if u.Role == "" {
		slog.Error("role is required or the role incorrect ", "error", u.Role)
		return ""
	}
	return u.Role
}

func (u *User) EmailWrapper() error {
	if validateEmail := u.ValidateEmail(); validateEmail == "" {
		slog.Error("failed to validate email", "email", u.Email)
		return ErrInvalidEmail
	}
	return nil
}

func (u *User) RoleWrapper() error {
	if validateRole := u.ValidateRole(); validateRole == "" {
		slog.Error("failed to validate role", "role", u.Role)
		return ErrInvalidRole
	}
	return nil
}
