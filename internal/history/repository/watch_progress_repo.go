package repository

import (
	"context"
	"fmt"

	"kinoza-back/internal/history/domain"
	"kinoza-back/pkg/querier"
)

type historyRepo struct {
	db querier.Querier
}

func NewWatchProgressRepository(db querier.Querier) domain.WatchProgressRepository {
	return &historyRepo{db: db}
}

func (r *historyRepo) Upsert(ctx context.Context, p *domain.WatchProgress) error {
	query := `INSERT INTO watch_progress(user_id, title_id, position_sec) VALUES($1, $2, $3) ON CONFLICT (user_id, title_id) DO UPDATE SET position_sec = EXCLUDED.position_sec, updated_at = NOW() RETURNING updated_at`

	err := querier.GetQuerier(ctx, r.db).GetContext(ctx, p, query, p.UserID, p.TitleID, p.PositionSec)
	if err != nil {
		return fmt.Errorf("history repo upsert: %w", err)
	}
	return nil
}

func (r *historyRepo) ListInProgress(ctx context.Context, userID string) ([]domain.WatchProgress, error) {
	query := `SELECT user_id, title_id, position_sec, updated_at FROM watch_progress WHERE user_id = $1 ORDER BY updated_at DESC`

	items := []domain.WatchProgress{}
	err := querier.GetQuerier(ctx, r.db).SelectContext(ctx, &items, query, userID)
	if err != nil {
		return nil, fmt.Errorf("history repo list in progress: %w", err)
	}
	return items, nil
}
