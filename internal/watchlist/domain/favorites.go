package domain

import (
	"context"
	"kinoza-back/pkg/querier"
	"time"
)

type WatchlistItem struct {
	ID      int       `db:"id"`
	UserID  string    `db:"user_id"`
	TitleID int       `db:"title_id"`
	Status  string    `db:"status"`
	AddedAt time.Time `db:"added_at"`
}

type WatchlistRepository interface {
	Create(ctx context.Context, q querier.Querier, item *WatchlistItem) error
	Update(ctx context.Context, q querier.Querier, item *WatchlistItem) error
	Delete(ctx context.Context, q querier.Querier, id int) error
	List(ctx context.Context, q querier.Querier, userID string) ([]WatchlistItem, error)
}
