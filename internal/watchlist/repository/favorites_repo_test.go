//go:build integration

package repository

import (
	"context"
	"errors"
	"kinoza-back/internal/watchlist/domain"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func setupDB(t *testing.T) *sqlx.DB {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL is not set")
	}

	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// createUser - создает юзера для теста, после теста удаляет его
// (его записи watchlist удалятся каскадом)
func createUser(t *testing.T, db *sqlx.DB, email string) string {
	var id string
	err := db.Get(&id, `INSERT INTO users(name, email, password_hash, role) VALUES('test', $1, 'hash', 'viewer') RETURNING id`, email)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	t.Cleanup(func() { db.Exec(`DELETE FROM users WHERE id = $1`, id) })
	return id
}

func TestRepoCreate(t *testing.T) {
	db := setupDB(t)
	repo := NewWatchlistRepository(db)
	ctx := context.Background()
	userID := createUser(t, db, "create@test.com")

	t.Run("success", func(t *testing.T) {
		item := &domain.WatchlistItem{UserID: userID, TitleID: 1, Status: domain.StatusPlanned}

		err := repo.Create(ctx, item)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if item.ID == 0 {
			t.Error("expected id to be set")
		}
		if item.AddedAt.IsZero() {
			t.Error("expected added_at to be set")
		}
	})

	t.Run("already exists", func(t *testing.T) {
		item := &domain.WatchlistItem{UserID: userID, TitleID: 1, Status: domain.StatusWatched}

		err := repo.Create(ctx, item)

		if !errors.Is(err, domain.ErrAlreadyExists) {
			t.Errorf("expected already exists error, got %v", err)
		}
	})
}

func TestRepoUpdate(t *testing.T) {
	db := setupDB(t)
	repo := NewWatchlistRepository(db)
	ctx := context.Background()
	ownerID := createUser(t, db, "update-owner@test.com")
	otherID := createUser(t, db, "update-other@test.com")

	item := &domain.WatchlistItem{UserID: ownerID, TitleID: 1, Status: domain.StatusPlanned}
	if err := repo.Create(ctx, item); err != nil {
		t.Fatalf("failed to create item: %v", err)
	}

	t.Run("other user", func(t *testing.T) {
		err := repo.Update(ctx, &domain.WatchlistItem{ID: item.ID, UserID: otherID, Status: domain.StatusWatched})

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected not found error, got %v", err)
		}
		var status string
		db.Get(&status, `SELECT status FROM watchlist_items WHERE id = $1`, item.ID)
		if status != domain.StatusPlanned {
			t.Errorf("status should not change, got %q", status)
		}
	})

	t.Run("not exists", func(t *testing.T) {
		err := repo.Update(ctx, &domain.WatchlistItem{ID: -1, UserID: ownerID, Status: domain.StatusWatched})

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected not found error, got %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		err := repo.Update(ctx, &domain.WatchlistItem{ID: item.ID, UserID: ownerID, Status: domain.StatusWatched})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		var status string
		db.Get(&status, `SELECT status FROM watchlist_items WHERE id = $1`, item.ID)
		if status != domain.StatusWatched {
			t.Errorf("expected status %q, got %q", domain.StatusWatched, status)
		}
	})
}

func TestRepoDelete(t *testing.T) {
	db := setupDB(t)
	repo := NewWatchlistRepository(db)
	ctx := context.Background()
	ownerID := createUser(t, db, "delete-owner@test.com")
	otherID := createUser(t, db, "delete-other@test.com")

	item := &domain.WatchlistItem{UserID: ownerID, TitleID: 1, Status: domain.StatusPlanned}
	if err := repo.Create(ctx, item); err != nil {
		t.Fatalf("failed to create item: %v", err)
	}

	t.Run("other user", func(t *testing.T) {
		err := repo.Delete(ctx, item.ID, otherID)

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected not found error, got %v", err)
		}

		var count int
		db.Get(&count, `SELECT COUNT(*) FROM watchlist_items WHERE id = $1`, item.ID)
		if count != 1 {
			t.Error("item should not be deleted by other user")
		}
	})

	t.Run("success", func(t *testing.T) {
		err := repo.Delete(ctx, item.ID, ownerID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("already deleted", func(t *testing.T) {
		err := repo.Delete(ctx, item.ID, ownerID)

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected not found error, got %v", err)
		}
	})
}

func TestRepoList(t *testing.T) {
	db := setupDB(t)
	repo := NewWatchlistRepository(db)
	ctx := context.Background()
	ownerID := createUser(t, db, "list-owner@test.com")
	otherID := createUser(t, db, "list-other@test.com")

	repo.Create(ctx, &domain.WatchlistItem{UserID: ownerID, TitleID: 1, Status: domain.StatusPlanned})
	repo.Create(ctx, &domain.WatchlistItem{UserID: ownerID, TitleID: 2, Status: domain.StatusWatching})
	repo.Create(ctx, &domain.WatchlistItem{UserID: otherID, TitleID: 3, Status: domain.StatusWatched})

	t.Run("only own items", func(t *testing.T) {
		items, err := repo.List(ctx, ownerID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(items))
		}
		for _, item := range items {
			if item.UserID != ownerID {
				t.Errorf("got item of other user: %+v", item)
			}
		}
	})

	t.Run("empty list", func(t *testing.T) {
		emptyID := createUser(t, db, "list-empty@test.com")

		items, err := repo.List(ctx, emptyID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if items == nil || len(items) != 0 {
			t.Errorf("expected empty list, got %v", items)
		}
	})
}
