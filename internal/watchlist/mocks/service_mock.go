package mocks

import (
	"context"

	"kinoza-back/internal/watchlist/domain"

	"github.com/stretchr/testify/mock"
)

type WatchlistServiceMock struct {
	mock.Mock
}

var _ domain.WatchlistService = (*WatchlistServiceMock)(nil)

func (m *WatchlistServiceMock) Create(ctx context.Context, item *domain.WatchlistItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *WatchlistServiceMock) Update(ctx context.Context, item *domain.WatchlistItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *WatchlistServiceMock) Delete(ctx context.Context, titleID int, userID string) error {
	args := m.Called(ctx, titleID, userID)
	return args.Error(0)
}

func (m *WatchlistServiceMock) List(ctx context.Context, userID string) ([]domain.WatchlistItem, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.WatchlistItem), args.Error(1)
}
