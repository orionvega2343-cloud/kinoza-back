package domain

import (
	"context"
	"time"
)

const (
	StatusPlanned  = "planned"
	StatusWatching = "watching"
	StatusWatched  = "watched"
)

func IsValidStatus(status string) bool {
	switch status {
	case StatusPlanned, StatusWatching, StatusWatched:
		return true
	default:
		return false
	}
}

type WatchlistItem struct {
	ID      int       `db:"id"`
	UserID  string    `db:"user_id"`
	TitleID int       `db:"title_id"`
	Status  string    `db:"status"`
	AddedAt time.Time `db:"added_at"`
}

type WatchlistRepository interface {
	Create(ctx context.Context, item *WatchlistItem) error
	Update(ctx context.Context, item *WatchlistItem) error
	Delete(ctx context.Context, id int, userID string) error
	List(ctx context.Context, userID string) ([]WatchlistItem, error)
}

type WatchlistService interface {
	Create(ctx context.Context, item *WatchlistItem) error
	Update(ctx context.Context, item *WatchlistItem) error
	Delete(ctx context.Context, id int, userID string) error
	List(ctx context.Context, userID string) ([]WatchlistItem, error)
}
