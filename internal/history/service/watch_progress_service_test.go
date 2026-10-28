package service

import (
	"context"
	"testing"

	"kinoza-back/internal/history/domain"
	"kinoza-back/internal/history/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSave(t *testing.T) {
	tests := []struct {
		name     string
		position int
		repoErr  error
		wantErr  error
		wantCall bool
	}{
		{name: "valid position", position: 1200, wantErr: nil, wantCall: true},
		{name: "zero position", position: 0, wantErr: nil, wantCall: true},
		{name: "negative position", position: -5, wantErr: domain.ErrInvalidPosition, wantCall: false},
		{name: "repo error", position: 1200, repoErr: domain.ErrDatabase, wantErr: domain.ErrDatabase, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.WatchProgressRepoMock)
			if tt.wantCall {
				repo.On("Upsert", mock.Anything, mock.Anything).Return(tt.repoErr)
			}
			svc := NewWatchProgressService(repo)

			err := svc.Save(context.Background(), &domain.WatchProgress{PositionSec: tt.position})

			assert.ErrorIs(t, err, tt.wantErr)
			repo.AssertExpectations(t)
		})
	}
}

func TestListInProgress(t *testing.T) {
	items := []domain.WatchProgress{{UserID: "user-1", TitleID: 42, PositionSec: 1200}}

	tests := []struct {
		name      string
		repoItems []domain.WatchProgress
		repoErr   error
		wantItems []domain.WatchProgress
		wantErr   error
	}{
		{name: "success", repoItems: items, wantItems: items, wantErr: nil},
		{name: "repo error", repoErr: domain.ErrDatabase, wantErr: domain.ErrDatabase},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.WatchProgressRepoMock)
			repo.On("ListInProgress", mock.Anything, "user-1").Return(tt.repoItems, tt.repoErr)
			svc := NewWatchProgressService(repo)

			got, err := svc.ListInProgress(context.Background(), "user-1")

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantItems, got)
			repo.AssertExpectations(t)
		})
	}
}
