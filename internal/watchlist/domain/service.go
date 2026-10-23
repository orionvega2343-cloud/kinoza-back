package domain

import "context"

type WatchlistService interface {
	Create(ctx context.Context, item *WatchlistItem) error
	Update(ctx context.Context, item *WatchlistItem) error
	Delete(ctx context.Context, titleID int, userID string) error
	List(ctx context.Context, userID string) ([]WatchlistItem, error)
}
