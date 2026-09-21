package service

import (
	"context"

	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/internal/watchlist/repository"
)

type WatchlistService struct {
	repo repository.WatchlistRepository
}

func NewWatchlistService(repo repository.WatchlistRepository) *WatchlistService {
	return &WatchlistService{
		repo: repo,
	}
}

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
