package service

import (
	"context"

	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/internal/watchlist/repository"
)

type WatchlistService struct {
	repo repository.WatchlistRepository
}

// Constructor
func NewWatchlistService(repo repository.WatchlistRepository) *WatchlistService {
	return &WatchlistService{
		repo: repo,
	}
}

func (s *WatchlistService) Create(ctx context.Context, item *domain.WatchlistItem) error {
	return s.repo.Create(ctx, item)
}
