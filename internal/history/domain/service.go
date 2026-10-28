package domain

import "context"

type WatchProgressService interface {
	Save(ctx context.Context, p *WatchProgress) error
	ListInProgress(ctx context.Context, userID string) ([]WatchProgress, error)
}
