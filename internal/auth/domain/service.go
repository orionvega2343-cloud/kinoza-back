package domain

import "context"

type UserService interface {
	Register(ctx context.Context, u *User) (*User, error)
	Login(ctx context.Context, email, password string) (access string, refresh string, err error)
	GetUserById(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, u *User, name, email string) error
	Refresh(ctx context.Context, rawToken string) (access, refresh string, err error)
	GetTokensList(ctx context.Context, userId string) ([]*RefreshToken, error)
	RevokeToken(ctx context.Context, id string) error
	RevokeAllTokens(ctx context.Context, userId string) error
}
