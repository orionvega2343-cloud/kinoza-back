package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/pkg/querier"
)

type watchlistRepo struct {
	db querier.Querier
}

func NewWatchlistRepository(db querier.Querier) domain.WatchlistRepository {
	return &watchlistRepo{db: db}
}

func (r *watchlistRepo) Create(ctx context.Context, item *domain.WatchlistItem) error {
	query := `INSERT INTO watchlist_items(user_id, title_id, status) VALUES($1, $2, $3) ON CONFLICT (user_id, title_id) DO NOTHING RETURNING id, added_at`

	err := querier.GetQuerier(ctx, r.db).GetContext(ctx, item, query, item.UserID, item.TitleID, item.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("watchlist repo create: %w", err)
	}
	return nil
}

func (r *watchlistRepo) Update(ctx context.Context, item *domain.WatchlistItem) error {
	query := `UPDATE watchlist_items SET status = $1 WHERE id = $2`

	result, err := querier.GetQuerier(ctx, r.db).ExecContext(ctx, query, item.Status, item.ID)
	if err != nil {
		return fmt.Errorf("watchlist repo update: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("watchlist repo update: %w", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *watchlistRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM watchlist_items WHERE id = $1`

	result, err := querier.GetQuerier(ctx, r.db).ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("watchlist repo delete: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("watchlist repo delete: %w", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *watchlistRepo) List(ctx context.Context, userID string) ([]domain.WatchlistItem, error) {
	query := `SELECT id, user_id, title_id, status, added_at FROM watchlist_items WHERE user_id = $1 ORDER BY added_at DESC`

	items := []domain.WatchlistItem{}
	err := querier.GetQuerier(ctx, r.db).SelectContext(ctx, &items, query, userID)
	if err != nil {
		return nil, fmt.Errorf("watchlist repo list: %w", err)
	}
	return items, nil
}
