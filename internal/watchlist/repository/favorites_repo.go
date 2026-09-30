package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/pkg/querier"
	"kinoza-back/pkg/transaction"

	"github.com/jmoiron/sqlx"
)

type watchlistRepo struct {
	db *sqlx.DB
}

func NewWatchListRepository(db *sqlx.DB) domain.WatchlistRepository {
	return &watchlistRepo{
		db: db,
	}
}

func (r *watchlistRepo) getQuerier(ctx context.Context) querier.Querier {
	tx, ok := transaction.ExtractTx(ctx)
	if ok {
		return tx
	}
	return r.db
}

func (r *watchlistRepo) Create(ctx context.Context, item *domain.WatchlistItem) error {
	query := `INSERT INTO watchlist_items(user_id, title_id, status) VALUES($1, $2, $3) ON CONFLICT (user_id, title_id) DO NOTHING RETURNING id, added_at`

	err := r.getQuerier(ctx).GetContext(ctx, item, query, item.UserID, item.TitleID, item.Status)
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

	result, err := r.getQuerier(ctx).ExecContext(ctx, query, item.Status, item.ID)
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

	result, err := r.getQuerier(ctx).ExecContext(ctx, query, id)
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
	err := r.getQuerier(ctx).SelectContext(ctx, &items, query, userID)
	if err != nil {
		return nil, fmt.Errorf("watchlist repo list: %w", err)
	}
	return items, err
}
