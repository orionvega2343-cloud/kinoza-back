package service

import (
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
