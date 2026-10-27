package service

import (
	"context"
	"fmt"

	"kinoza-back/internal/history/domain"
)

type WatchProgressService struct {
	repo domain.WatchProgressRepository
}

func NewWatchProgressService(repo domain.WatchProgressRepository) *WatchProgressService {
	return &WatchProgressService{
		repo: repo,
	}
}

var _ domain.WatchProgressService = (*WatchProgressService)(nil)

func (s *WatchProgressService) Save(ctx context.Context, p *domain.WatchProgress) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("history service: Save, %w", err)
	}
	return s.repo.Upsert(ctx, p)
}

func (s *WatchProgressService) ListInProgress(ctx context.Context, userID string) ([]domain.WatchProgress, error) {
	return s.repo.ListInProgress(ctx, userID)
}
