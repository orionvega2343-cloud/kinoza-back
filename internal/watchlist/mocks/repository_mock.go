package mocks

import (
	"context"

	"kinoza-back/internal/watchlist/domain"

	"github.com/stretchr/testify/mock"
)

type WatchlistRepoMock struct {
	mock.Mock
}

var _ domain.WatchlistRepository = (*WatchlistRepoMock)(nil)

func (m *WatchlistRepoMock) Create(ctx context.Context, item *domain.WatchlistItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *WatchlistRepoMock) Update(ctx context.Context, item *domain.WatchlistItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *WatchlistRepoMock) Delete(ctx context.Context, titleID int, userID string) error {
	args := m.Called(ctx, titleID, userID)
	return args.Error(0)
}

func (m *WatchlistRepoMock) List(ctx context.Context, userID string) ([]domain.WatchlistItem, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.WatchlistItem), args.Error(1)
}
