package domain

import (
	"context"
	"time"
)

type WatchlistItem struct {
	ID      int    `db:"id"`
	UserID  string `db:"user_id"`
	TitleID int    `db:"title_id"`
	// Status - planed | watching | watched
	Status  string    `db:"status"`
	AddedAt time.Time `db:"added_at"`
}

type WatchlistRepository interface {
	Create(ctx context.Context, item *WatchlistItem) error
	Update(ctx context.Context, item *WatchlistItem) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, userID string) ([]WatchlistItem, error)
}

type WatchlistService interface {
	Create(ctx context.Context, item *WatchlistItem) error
	Update(ctx context.Context, item *WatchlistItem) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, userID string) ([]WatchlistItem, error)
}
