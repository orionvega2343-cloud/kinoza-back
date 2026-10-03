package domain

import "context"

type UserRepo interface {
	CreateUser(ctx context.Context, u *User) (*User, error)
	GetUserById(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, u *User, name, email string) error
}
