package repository

import (
	"context"
	"errors"
	"fmt"

	"kinoza-back/internal/friends/domain"
	"kinoza-back/pkg/querier"

	"github.com/lib/pq"
)

type friendshipRepo struct {
	db querier.Querier
}

func NewFriendshipRepository(db querier.Querier) domain.FriendshipRepository {
	return &friendshipRepo{db: db}
}

func (r *friendshipRepo) Create(ctx context.Context, f *domain.Friendship) error {
	query := `
	INSERT INTO friendships (user_id_a, user_id_b)
	VALUES ($1, $2)
	RETURNING id, created_at
	`
	err := querier.GetQuerier(ctx, r.db).GetContext(ctx, f, query, f.UserIDA, f.UserIDB)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqUniqueViolation {
			return domain.ErrAlreadyFriends
		}
		return fmt.Errorf("friendship repo create: %w", err)
	}
	return nil
}

func (r *friendshipRepo) Exists(ctx context.Context, f *domain.Friendship) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM friendships WHERE user_id_a = $1 AND user_id_b = $2)`

	var exists bool
	err := querier.GetQuerier(ctx, r.db).GetContext(ctx, &exists, query, f.UserIDA, f.UserIDB)
	if err != nil {
		return false, fmt.Errorf("friendship repo exists: %w", err)
	}
	return exists, nil
}

func (r *friendshipRepo) Delete(ctx context.Context, f *domain.Friendship) error {
	query := `DELETE FROM friendships WHERE user_id_a = $1 AND user_id_b = $2`

	result, err := querier.GetQuerier(ctx, r.db).ExecContext(ctx, query, f.UserIDA, f.UserIDB)
	if err != nil {
		return fmt.Errorf("friendship repo delete: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("friendship repo delete: %w", err)
	}

	if rows == 0 {
		return domain.ErrFriendshipNotFound
	}
	return nil
}

func (r *friendshipRepo) List(ctx context.Context, userID string) ([]domain.Friendship, error) {
	query := `
	SELECT id, user_id_a, user_id_b, created_at
	FROM friendships
	WHERE user_id_a = $1 OR user_id_b = $1
	ORDER BY created_at DESC
	`
	friendships := []domain.Friendship{}
	err := querier.GetQuerier(ctx, r.db).SelectContext(ctx, &friendships, query, userID)
	if err != nil {
		return nil, fmt.Errorf("friendship repo list: %w", err)
	}
	return friendships, nil
}
