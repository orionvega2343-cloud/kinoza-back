package mocks

import (
	"context"
	"kinoza-back/internal/auth/domain"

	"github.com/stretchr/testify/mock"
)

type TokenMock struct {
	mock.Mock
}

func (t *TokenMock) CreateToken(ctx context.Context, token *domain.RefreshToken) (*domain.RefreshToken, error) {
	args := t.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (t *TokenMock) GetTokenByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	args := t.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (t *TokenMock) GetTokensList(ctx context.Context, userId string) ([]*domain.RefreshToken, error) {
	args := t.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.RefreshToken), args.Error(1)
}

func (t *TokenMock) RevokeToken(ctx context.Context, id string) error {
	args := t.Called(ctx, id)
	return args.Error(0)
}

func (t *TokenMock) RevokeAllTokens(ctx context.Context, userId string) error {
	args := t.Called(ctx, userId)
	return args.Error(0)
}
