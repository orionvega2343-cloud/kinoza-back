package mocks

import (
	"context"

	"kinoza-back/internal/history/domain"

	"github.com/stretchr/testify/mock"
)

type WatchProgressRepoMock struct {
	mock.Mock
}

var _ domain.WatchProgressRepository = (*WatchProgressRepoMock)(nil)

func (m *WatchProgressRepoMock) Upsert(ctx context.Context, p *domain.WatchProgress) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *WatchProgressRepoMock) ListInProgress(ctx context.Context, userID string) ([]domain.WatchProgress, error) {
	args := m.Called(ctx, userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.WatchProgress), args.Error(1)
}
