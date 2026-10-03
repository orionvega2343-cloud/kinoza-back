package service

import (
	"context"

	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/pkg/transaction"
)

type WatchlistService struct {
	repo domain.WatchlistRepository
	tx   *transaction.Transactor
}

func NewWatchlistService(repo domain.WatchlistRepository, tx *transaction.Transactor) *WatchlistService {
	return &WatchlistService{
		repo: repo,
		tx:   tx,
	}
}

var _ domain.WatchlistService = (*WatchlistService)(nil)

func (s *WatchlistService) Create(ctx context.Context, item *domain.WatchlistItem) error {
	return s.repo.Create(ctx, item)
}

func (s *WatchlistService) Update(ctx context.Context, item *domain.WatchlistItem) error {
	return s.repo.Update(ctx, item)
}

func (s *WatchlistService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *WatchlistService) List(ctx context.Context, userID string) ([]domain.WatchlistItem, error) {
	return s.repo.List(ctx, userID)
}
