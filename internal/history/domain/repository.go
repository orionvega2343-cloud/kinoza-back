package domain

import "context"

type WatchProgressRepository interface {
	Upsert(ctx context.Context, p *WatchProgress) error
	Get(ctx context.Context, userID string, titleID int) (*WatchProgress, error)
	ListInProgress(ctx context.Context, userID string) ([]WatchProgress, error)
}
