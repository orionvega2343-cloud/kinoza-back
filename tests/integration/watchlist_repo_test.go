//go:build integration

package integration

import (
	"context"
	"testing"

	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/internal/watchlist/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatchlistRepo_Create(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	repo := repository.NewWatchlistRepository(db)
	userID := createUser(t)

	t.Run("creates new item", func(t *testing.T) {
		item := &domain.WatchlistItem{UserID: userID, TitleID: 1, Status: domain.StatusPlanned}

		require.NoError(t, repo.Create(ctx, item))
		assert.NotZero(t, item.ID)
		assert.False(t, item.AddedAt.IsZero())
	})

	t.Run("same title again is upsert, not duplicate", func(t *testing.T) {
		first := &domain.WatchlistItem{UserID: userID, TitleID: 2, Status: domain.StatusPlanned}
		require.NoError(t, repo.Create(ctx, first))

		second := &domain.WatchlistItem{UserID: userID, TitleID: 2, Status: domain.StatusWatching}
		require.NoError(t, repo.Create(ctx, second))
		assert.Equal(t, first.ID, second.ID)

		var count int
		var status string
		require.NoError(t, db.Get(&count, `SELECT COUNT(*) FROM watchlist_items WHERE user_id = $1 AND title_id = 2`, userID))
		require.NoError(t, db.Get(&status, `SELECT status FROM watchlist_items WHERE user_id = $1 AND title_id = 2`, userID))
		assert.Equal(t, 1, count)
		assert.Equal(t, domain.StatusWatching, status)
	})

	t.Run("invalid status rejected by CHECK", func(t *testing.T) {
		item := &domain.WatchlistItem{UserID: userID, TitleID: 3, Status: "bad"}

		assert.Error(t, repo.Create(ctx, item))
	})

	t.Run("unknown user rejected by FK", func(t *testing.T) {
		item := &domain.WatchlistItem{UserID: "00000000-0000-0000-0000-000000000000", TitleID: 1, Status: domain.StatusPlanned}

		assert.Error(t, repo.Create(ctx, item))
	})
}

func TestWatchlistRepo_Update(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	repo := repository.NewWatchlistRepository(db)
	userID := createUser(t)
	otherID := createUser(t)
	require.NoError(t, repo.Create(ctx, &domain.WatchlistItem{UserID: userID, TitleID: 1, Status: domain.StatusPlanned}))

	t.Run("updates status", func(t *testing.T) {
		require.NoError(t, repo.Update(ctx, &domain.WatchlistItem{UserID: userID, TitleID: 1, Status: domain.StatusWatched}))

		var status string
		require.NoError(t, db.Get(&status, `SELECT status FROM watchlist_items WHERE user_id = $1 AND title_id = 1`, userID))
		assert.Equal(t, domain.StatusWatched, status)
	})

	t.Run("missing item returns ErrNotFound", func(t *testing.T) {
		err := repo.Update(ctx, &domain.WatchlistItem{UserID: userID, TitleID: 999, Status: domain.StatusWatched})

		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("other user's item returns ErrNotFound", func(t *testing.T) {
		err := repo.Update(ctx, &domain.WatchlistItem{UserID: otherID, TitleID: 1, Status: domain.StatusPlanned})

		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestWatchlistRepo_Delete(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	repo := repository.NewWatchlistRepository(db)
	userID := createUser(t)
	otherID := createUser(t)
	require.NoError(t, repo.Create(ctx, &domain.WatchlistItem{UserID: userID, TitleID: 1, Status: domain.StatusPlanned}))

	t.Run("other user's item returns ErrNotFound and is kept", func(t *testing.T) {
		assert.ErrorIs(t, repo.Delete(ctx, 1, otherID), domain.ErrNotFound)

		items, err := repo.List(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, items, 1)
	})

	t.Run("deletes item", func(t *testing.T) {
		require.NoError(t, repo.Delete(ctx, 1, userID))

		items, err := repo.List(ctx, userID)
		require.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("second delete returns ErrNotFound", func(t *testing.T) {
		assert.ErrorIs(t, repo.Delete(ctx, 1, userID), domain.ErrNotFound)
	})
}

func TestWatchlistRepo_List(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	repo := repository.NewWatchlistRepository(db)
	userID := createUser(t)
	otherID := createUser(t)

	t.Run("empty list is not nil", func(t *testing.T) {
		items, err := repo.List(ctx, userID)

		require.NoError(t, err)
		assert.NotNil(t, items)
		assert.Empty(t, items)
	})

	t.Run("only own items, newest first", func(t *testing.T) {
		require.NoError(t, repo.Create(ctx, &domain.WatchlistItem{UserID: userID, TitleID: 1, Status: domain.StatusPlanned}))
		require.NoError(t, repo.Create(ctx, &domain.WatchlistItem{UserID: userID, TitleID: 2, Status: domain.StatusPlanned}))
		require.NoError(t, repo.Create(ctx, &domain.WatchlistItem{UserID: otherID, TitleID: 3, Status: domain.StatusPlanned}))

		items, err := repo.List(ctx, userID)
		require.NoError(t, err)
		require.Len(t, items, 2)
		assert.Equal(t, 2, items[0].TitleID)
		assert.Equal(t, 1, items[1].TitleID)
		for _, item := range items {
			assert.Equal(t, userID, item.UserID)
		}
	})
}
