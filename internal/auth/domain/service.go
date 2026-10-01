package domain

import "context"

type UserService interface {
	Register(ctx context.Context, u *User) (*User, error)
	Login(ctx context.Context, email string, password string) (access string, refresh string, err error)
	GetUserById(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, u *User, name, email string) error
}
