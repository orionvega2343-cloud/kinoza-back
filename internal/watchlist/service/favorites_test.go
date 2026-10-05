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

func TestCreate_InvalidStatus(t *testing.T) {
	repo := &fakeRepo{}
	src := NewWatchlistService(repo, nil)

	err := src.Create(context.Background(), &domain.WatchlistItem{Status: "lol"})

	if !errors.Is(err, domain.ErrInvalidStatus) {
		t.Errorf("expected invalid status error, got %v", err)
	}

	if repo.call {
		t.Error("repo should not be called with invalid status")
	}
}

func TestCreate_ValidStatus(t *testing.T) {
	repo := &fakeRepo{}
	src := NewWatchlistService(repo, nil)

	err := src.Create(context.Background(), &domain.WatchlistItem{Status: domain.StatusPlanned})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !repo.call {
		t.Error("repo should be called with valid status")
	}
}

func TestCreate_RepoError(t *testing.T) {
	repo := &fakeRepo{err: domain.ErrAlreadyExists}
	src := NewWatchlistService(repo, nil)

	err := src.Create(context.Background(), &domain.WatchlistItem{Status: domain.StatusPlanned})

	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Errorf("expected already exists error, got %v", err)
	}
}

func TestUpdate_InvalidStatus(t *testing.T) {
	repo := &fakeRepo{}
	src := NewWatchlistService(repo, nil)

	err := src.Update(context.Background(), &domain.WatchlistItem{Status: "lol"})

	if !errors.Is(err, domain.ErrInvalidStatus) {
		t.Errorf("expected invalid status error, got %v", err)
	}

	if repo.call {
		t.Error("repo should not be called with invalid status")
	}
}

func TestUpdate_ValidStatus(t *testing.T) {
	repo := &fakeRepo{}
	src := NewWatchlistService(repo, nil)

	err := src.Update(context.Background(), &domain.WatchlistItem{Status: domain.StatusWatched})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !repo.call {
		t.Error("repo should be called with valid status")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	repo := &fakeRepo{err: domain.ErrNotFound}
	src := NewWatchlistService(repo, nil)

	err := src.Update(context.Background(), &domain.WatchlistItem{Status: domain.StatusWatched})

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	repo := &fakeRepo{err: domain.ErrNotFound}
	src := NewWatchlistService(repo, nil)

	err := src.Delete(context.Background(), 1, "user-1")

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected not found error, got %v", err)
	}

	if !repo.call {
		t.Error("repo should be called")
	}
}

func TestList(t *testing.T) {
	repo := &fakeRepo{}
	src := NewWatchlistService(repo, nil)

	_, err := src.List(context.Background(), "user-1")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !repo.call {
		t.Error("repo should be called")
	}
}
