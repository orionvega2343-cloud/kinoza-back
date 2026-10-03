package repository

import (
	"context"
	"kinoza-back/internal/auth/domain"
	"kinoza-back/pkg/logger"
	"kinoza-back/pkg/querier"
	"log/slog"
)

var _ domain.RefreshTokens = (*RefreshTokenImpl)(nil)

type RefreshTokenImpl struct {
	q querier.Querier
}

func NewRefreshTokenRepo(q querier.Querier) *RefreshTokenImpl {
	return &RefreshTokenImpl{q: q}
}

func (t *RefreshTokenImpl) CreateToken(ctx context.Context, token *domain.RefreshToken) (*domain.RefreshToken, error) {
	query := `INSERT INTO refresh_tokens(user_id, token_hash, expires_at) VALUES($1, $2, $3) RETURNING id, user_id, token_hash, expires_at, revoked_at, created_at`
	if err := querier.GetQuerier(ctx, t.q).GetContext(ctx, token, query, token.UserId, token.TokenHash, token.ExpiresAt); err != nil {
		return nil, logger.LogErr("failed to insert refresh token", err)
	}
	return token, nil
}

func (t *RefreshTokenImpl) GetTokensList(ctx context.Context, userId string) ([]*domain.RefreshToken, error) {
	token := []*domain.RefreshToken{}
	query := `SELECT id, user_id, token_hash, expires_at, revoked_at, created_at FROM refresh_tokens WHERE user_id = $1`
	if err := querier.GetQuerier(ctx, t.q).SelectContext(ctx, &token, query, userId); err != nil {
		return nil, logger.LogErr("failed to fetch refresh tokens", err)
	}
	return token, nil
}

func (t *RefreshTokenImpl) GetTokenByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	token := &domain.RefreshToken{}
	query := `SELECT id, user_id, token_hash, expires_at, revoked_at, created_at FROM refresh_tokens WHERE token_hash = $1`
	if err := querier.GetQuerier(ctx, t.q).GetContext(ctx, token, query, tokenHash); err != nil {
		return nil, logger.LogErr("failed to fetch refresh token", err)
	}
	return token, nil
}

func (t *RefreshTokenImpl) RevokeToken(ctx context.Context, id string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id = $1`
	if _, err := querier.GetQuerier(ctx, t.q).ExecContext(ctx, query, id); err != nil {
		return logger.LogErr("failed to revoke refresh token", err)
	}
	return nil
}

func (t *RefreshTokenImpl) RevokeAllTokens(ctx context.Context, userId string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`
	if _, err := querier.GetQuerier(ctx, t.q).ExecContext(ctx, query, userId); err != nil {
		slog.Error("failed to revoke all refresh tokens", "error", err)
		return logger.LogErr("failed to revoke all refresh tokens", err)
	}
	return nil
}
