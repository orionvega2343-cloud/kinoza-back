package service

import (
	"context"
	"testing"

	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/internal/watchlist/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		repoErr  error
		wantErr  error
		wantCall bool
	}{
		{name: "status planned", status: domain.StatusPlanned, wantErr: nil, wantCall: true},
		{name: "status watching", status: domain.StatusWatching, wantErr: nil, wantCall: true},
		{name: "status watched", status: domain.StatusWatched, wantErr: nil, wantCall: true},
		{name: "invalid status", status: "lol", wantErr: domain.ErrInvalidStatus, wantCall: false},
		{name: "empty status", status: "", wantErr: domain.ErrInvalidStatus, wantCall: false},
		{name: "uppercase status", status: "Planned", wantErr: domain.ErrInvalidStatus, wantCall: false},
		{name: "repo error", status: domain.StatusPlanned, repoErr: domain.ErrDB, wantErr: domain.ErrDB, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.WatchlistRepoMock)
			if tt.wantCall {
				repo.On("Create", mock.Anything, mock.Anything).Return(tt.repoErr)
			}
			svc := NewWatchlistService(repo, nil)

			err := svc.Create(context.Background(), &domain.WatchlistItem{Status: tt.status})

			assert.ErrorIs(t, err, tt.wantErr)
			repo.AssertExpectations(t)
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		repoErr  error
		wantErr  error
		wantCall bool
	}{
		{name: "status planned", status: domain.StatusPlanned, wantErr: nil, wantCall: true},
		{name: "status watching", status: domain.StatusWatching, wantErr: nil, wantCall: true},
		{name: "status watched", status: domain.StatusWatched, wantErr: nil, wantCall: true},
		{name: "invalid status", status: "lol", wantErr: domain.ErrInvalidStatus, wantCall: false},
		{name: "empty status", status: "", wantErr: domain.ErrInvalidStatus, wantCall: false},
		{name: "uppercase status", status: "Planned", wantErr: domain.ErrInvalidStatus, wantCall: false},
		{name: "not found", status: domain.StatusWatched, repoErr: domain.ErrNotFound, wantErr: domain.ErrNotFound, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.WatchlistRepoMock)
			if tt.wantCall {
				repo.On("Update", mock.Anything, mock.Anything).Return(tt.repoErr)
			}
			svc := NewWatchlistService(repo, nil)

			err := svc.Update(context.Background(), &domain.WatchlistItem{ID: 1, Status: tt.status})

			assert.ErrorIs(t, err, tt.wantErr)
			repo.AssertExpectations(t)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		repoErr error
		wantErr error
	}{
		{name: "success", wantErr: nil},
		{name: "not found", repoErr: domain.ErrNotFound, wantErr: domain.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.WatchlistRepoMock)
			repo.On("Delete", mock.Anything, 1, "user-1").Return(tt.repoErr)
			svc := NewWatchlistService(repo, nil)

			err := svc.Delete(context.Background(), 1, "user-1")

			assert.ErrorIs(t, err, tt.wantErr)
			repo.AssertExpectations(t)
		})
	}
}

func TestList(t *testing.T) {
	items := []domain.WatchlistItem{{ID: 1, UserID: "user-1", TitleID: 42, Status: domain.StatusPlanned}}

	tests := []struct {
		name      string
		repoItems []domain.WatchlistItem
		repoErr   error
		wantItems []domain.WatchlistItem
		wantErr   error
	}{
		{name: "success", repoItems: items, wantItems: items, wantErr: nil},
		{name: "repo error", repoErr: domain.ErrDB, wantErr: domain.ErrDB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.WatchlistRepoMock)
			repo.On("List", mock.Anything, "user-1").Return(tt.repoItems, tt.repoErr)
			svc := NewWatchlistService(repo, nil)

			got, err := svc.List(context.Background(), "user-1")

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantItems, got)
			repo.AssertExpectations(t)
		})
	}
}
