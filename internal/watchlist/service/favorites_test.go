package service

import (
	"context"
	"errors"
	"testing"

	"kinoza-back/internal/watchlist/domain"
)

type fakeRepo struct {
	call bool
	err  error
}

func (f *fakeRepo) Create(ctx context.Context, item *domain.WatchlistItem) error {
	f.call = true
	return f.err
}

func (f *fakeRepo) Update(ctx context.Context, item *domain.WatchlistItem) error {
	f.call = true
	return f.err
}

func (f *fakeRepo) Delete(ctx context.Context, id int, userID string) error {
	f.call = true
	return f.err
}

func (f *fakeRepo) List(ctx context.Context, userID string) ([]domain.WatchlistItem, error) {
	f.call = true
	return nil, f.err
}

var errDB = errors.New("db error")

func TestCreate(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		repoErr  error
		wantErr  error
		wantCall bool
	}{
		{name: "valid status", status: domain.StatusPlanned, wantErr: nil, wantCall: true},
		{name: "invalid status", status: "lol", wantErr: domain.ErrInvalidStatus, wantCall: false},
		{name: "empty status", status: "", wantErr: domain.ErrInvalidStatus, wantCall: false},
		{name: "already exists", status: domain.StatusPlanned, repoErr: domain.ErrAlreadyExists, wantErr: domain.ErrAlreadyExists, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{err: tt.repoErr}
			src := NewWatchlistService(repo, nil)

			err := src.Create(context.Background(), &domain.WatchlistItem{Status: tt.status})

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if repo.call != tt.wantCall {
				t.Errorf("expected repo call %v, got %v", tt.wantCall, repo.call)
			}
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
		{name: "valid status", status: domain.StatusWatched, wantErr: nil, wantCall: true},
		{name: "invalid status", status: "lol", wantErr: domain.ErrInvalidStatus, wantCall: false},
		{name: "not found", status: domain.StatusWatched, repoErr: domain.ErrNotFound, wantErr: domain.ErrNotFound, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{err: tt.repoErr}
			src := NewWatchlistService(repo, nil)

			err := src.Update(context.Background(), &domain.WatchlistItem{ID: 1, Status: tt.status})

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if repo.call != tt.wantCall {
				t.Errorf("expected repo call %v, got %v", tt.wantCall, repo.call)
			}
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
			repo := &fakeRepo{err: tt.repoErr}
			src := NewWatchlistService(repo, nil)

			err := src.Delete(context.Background(), 1, "user-1")

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if !repo.call {
				t.Error("repo should be called")
			}
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name    string
		repoErr error
		wantErr error
	}{
		{name: "success", wantErr: nil},
		{name: "repo error", repoErr: errDB, wantErr: errDB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{err: tt.repoErr}
			src := NewWatchlistService(repo, nil)

			_, err := src.List(context.Background(), "user-1")

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if !repo.call {
				t.Error("repo should be called")
			}
		})
	}
}
