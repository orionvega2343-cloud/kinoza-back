package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"kinoza-back/internal/friends/domain"
	"kinoza-back/pkg/querier"

	"github.com/lib/pq"
)

const (
	pqUniqueViolation     = "23505"
	pqForeignKeyViolation = "23503"
)

type friendRequestRepo struct {
	db querier.Querier
}

func NewFriendRequestRepository(db querier.Querier) domain.FriendRequestRepository {
	return &friendRequestRepo{db: db}
}

func (r *friendRequestRepo) Create(ctx context.Context, req *domain.FriendRequest) error {
	query := `
	INSERT INTO friend_requests (from_user_id, to_user_id)
	VALUES ($1, $2)
	RETURNING id, status, created_at
	`
	err := querier.GetQuerier(ctx, r.db).GetContext(ctx, req, query, req.FromUserID, req.ToUserID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case pqUniqueViolation:
				return domain.ErrRequestAlreadyExists
			case pqForeignKeyViolation:
				return domain.ErrUserNotFound
			}
		}
		return fmt.Errorf("friend request repo create: %w", err)
	}
	return nil
}

func (r *friendRequestRepo) GetByIDForUpdate(ctx context.Context, id int) (*domain.FriendRequest, error) {
	query := `
	SELECT id, from_user_id, to_user_id, status, created_at
	FROM friend_requests
	WHERE id = $1
	FOR UPDATE
	`
	var req domain.FriendRequest
	err := querier.GetQuerier(ctx, r.db).GetContext(ctx, &req, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRequestNotFound
		}
		return nil, fmt.Errorf("friend request repo get by id: %w", err)
	}
	return &req, nil
}

func (r *friendRequestRepo) UpdateStatus(ctx context.Context, req *domain.FriendRequest) error {
	query := `UPDATE friend_requests SET status = $1 WHERE id = $2`

	result, err := querier.GetQuerier(ctx, r.db).ExecContext(ctx, query, req.Status, req.ID)
	if err != nil {
		return fmt.Errorf("friend request repo update status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("friend request repo update status: %w", err)
	}

	if rows == 0 {
		return domain.ErrRequestNotFound
	}
	return nil
}

func (r *friendRequestRepo) ListIncoming(ctx context.Context, userID string) ([]domain.FriendRequest, error) {
	query := `
	SELECT id, from_user_id, to_user_id, status, created_at
	FROM friend_requests
	WHERE to_user_id = $1 AND status = 'pending'
	ORDER BY created_at DESC
	`
	requests := []domain.FriendRequest{}
	err := querier.GetQuerier(ctx, r.db).SelectContext(ctx, &requests, query, userID)
	if err != nil {
		return nil, fmt.Errorf("friend request repo list incoming: %w", err)
	}
	return requests, nil
}

func (r *friendRequestRepo) ListOutgoing(ctx context.Context, userID string) ([]domain.FriendRequest, error) {
	query := `
	SELECT id, from_user_id, to_user_id, status, created_at
	FROM friend_requests
	WHERE from_user_id = $1 AND status = 'pending'
	ORDER BY created_at DESC
	`
	requests := []domain.FriendRequest{}
	err := querier.GetQuerier(ctx, r.db).SelectContext(ctx, &requests, query, userID)
	if err != nil {
		return nil, fmt.Errorf("friend request repo list outgoing: %w", err)
	}
	return requests, nil
}
