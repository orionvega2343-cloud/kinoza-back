package domain

import "context"

type RefreshTokens interface {
	CreateToken(ctx context.Context, token *RefreshToken) (*RefreshToken, error)
	GetTokenByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	GetTokensList(ctx context.Context, userId string) ([]*RefreshToken, error)
	RevokeToken(ctx context.Context, id string) error
}
