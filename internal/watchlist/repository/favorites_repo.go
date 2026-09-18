package repository

import (
	"context"

	"kinoza-back/internal/watchlist/domain"
)

type WatchlistRepository interface {
	Create(ctx context.Context, item *domain.WatchlistItem) error
	Update(ctx context.Context, item *domain.WatchlistItem) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, userID uint) ([]domain.WatchlistItem, error)
}
