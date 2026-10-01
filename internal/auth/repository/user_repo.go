package repository

import (
	"context"
	"kinoza-back/internal/auth/domain"
	"kinoza-back/pkg/logger"
	"kinoza-back/pkg/querier"
)

var _ domain.UserRepo = (*UserRepoImpl)(nil)

type UserRepoImpl struct {
	q querier.Querier
}

func NewUserRepoImpl(q querier.Querier) *UserRepoImpl {
	return &UserRepoImpl{q: q}
}

func (r *UserRepoImpl) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id`
	if err := querier.GetQuerier(ctx, r.q).GetContext(ctx, &u.Id, query, u.Email, u.PasswordHash, u.Role); err != nil {
		return nil, logger.LogErr("failed to insert user", err)
	}
	return u, nil
}

func (r *UserRepoImpl) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE id = $1`
	if err := querier.GetQuerier(ctx, r.q).GetContext(ctx, u, query, id); err != nil {
		return nil, logger.LogErr("failed to get user by id", err)
	}
	return u, nil
}

func (r *UserRepoImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`
	if err := querier.GetQuerier(ctx, r.q).GetContext(ctx, u, query, email); err != nil {
		return nil, logger.LogErr("failed to get user by email", err)
	}
	return u, nil
}

func (r *UserRepoImpl) UpdateUser(ctx context.Context, u *domain.User, name, email string) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	if _, err := querier.GetQuerier(ctx, r.q).ExecContext(ctx, query, name, email, u.Id); err != nil {
		return logger.LogErr("failed to update user", err)
	}
	return nil
}
