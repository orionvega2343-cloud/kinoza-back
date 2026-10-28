package domain

import "context"

type WatchProgressRepository interface {
	Upsert(ctx context.Context, p *WatchProgress) error
	ListInProgress(ctx context.Context, userID string) ([]WatchProgress, error)
}
