package service

import (
	"context"

	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/pkg/querier"
	"kinoza-back/pkg/transaction"
)

type WatchlistService struct {
	repo domain.WatchlistRepository
	tx   *transaction.Transactor
	db   querier.Querier
}

func NewWatchlistService(repo domain.WatchlistRepository, tx *transaction.Transactor, db querier.Querier) *WatchlistService {
	return &WatchlistService{
		repo: repo,
		tx:   tx,
		db:   db,
	}
}

func (s *WatchlistService) Create(ctx context.Context, item *domain.WatchlistItem) error {
	return s.repo.Create(ctx, s.db, item)
}

func (s *WatchlistService) Update(ctx context.Context, item *domain.WatchlistItem) error {
	return s.repo.Update(ctx, s.db, item)
}

func (s *WatchlistService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, s.db, id)
}

func (s *WatchlistService) List(ctx context.Context, userID string) ([]domain.WatchlistItem, error) {
	return s.repo.List(ctx, s.db, userID)
}
