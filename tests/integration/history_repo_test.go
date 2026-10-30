//go:build integration

package integration

import (
	"context"
	"testing"

	"kinoza-back/internal/history/domain"
	"kinoza-back/internal/history/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistoryRepo_Upsert(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	repo := repository.NewWatchProgressRepository(db)
	userID := createUser(t)

	t.Run("creates new record", func(t *testing.T) {
		p := &domain.WatchProgress{UserID: userID, TitleID: 1, PositionSec: 100}

		require.NoError(t, repo.Upsert(ctx, p))
		assert.False(t, p.UpdatedAt.IsZero())

		items, err := repo.ListInProgress(ctx, userID)
		require.NoError(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, 100, items[0].PositionSec)
	})

	t.Run("second upsert updates position, no duplicate", func(t *testing.T) {
		first := &domain.WatchProgress{UserID: userID, TitleID: 2, PositionSec: 100}
		require.NoError(t, repo.Upsert(ctx, first))

		second := &domain.WatchProgress{UserID: userID, TitleID: 2, PositionSec: 500}
		require.NoError(t, repo.Upsert(ctx, second))
		assert.False(t, second.UpdatedAt.Before(first.UpdatedAt))

		var count, position int
		require.NoError(t, db.Get(&count, `SELECT COUNT(*) FROM watch_progress WHERE user_id = $1 AND title_id = 2`, userID))
		require.NoError(t, db.Get(&position, `SELECT position_sec FROM watch_progress WHERE user_id = $1 AND title_id = 2`, userID))
		assert.Equal(t, 1, count)
		assert.Equal(t, 500, position)
	})

	t.Run("negative position rejected by CHECK", func(t *testing.T) {
		p := &domain.WatchProgress{UserID: userID, TitleID: 3, PositionSec: -1}

		assert.Error(t, repo.Upsert(ctx, p))
	})

	t.Run("unknown user rejected by FK", func(t *testing.T) {
		p := &domain.WatchProgress{UserID: "00000000-0000-0000-0000-000000000000", TitleID: 1, PositionSec: 10}

		assert.Error(t, repo.Upsert(ctx, p))
	})
}

func TestHistoryRepo_ListInProgress(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	repo := repository.NewWatchProgressRepository(db)
	userID := createUser(t)
	otherID := createUser(t)

	t.Run("empty list is not nil", func(t *testing.T) {
		items, err := repo.ListInProgress(ctx, userID)

		require.NoError(t, err)
		assert.NotNil(t, items)
		assert.Empty(t, items)
	})

	t.Run("only own records, last watched first", func(t *testing.T) {
		require.NoError(t, repo.Upsert(ctx, &domain.WatchProgress{UserID: userID, TitleID: 1, PositionSec: 10}))
		require.NoError(t, repo.Upsert(ctx, &domain.WatchProgress{UserID: userID, TitleID: 2, PositionSec: 20}))
		require.NoError(t, repo.Upsert(ctx, &domain.WatchProgress{UserID: otherID, TitleID: 3, PositionSec: 30}))
		// тайтл 1 пересмотрели последним, он должен подняться наверх
		require.NoError(t, repo.Upsert(ctx, &domain.WatchProgress{UserID: userID, TitleID: 1, PositionSec: 15}))

		items, err := repo.ListInProgress(ctx, userID)
		require.NoError(t, err)
		require.Len(t, items, 2)
		assert.Equal(t, 1, items[0].TitleID)
		assert.Equal(t, 15, items[0].PositionSec)
		assert.Equal(t, 2, items[1].TitleID)
		for _, item := range items {
			assert.Equal(t, userID, item.UserID)
		}
	})

	t.Run("user deleted, progress deleted by cascade", func(t *testing.T) {
		_, err := db.Exec(`DELETE FROM users WHERE id = $1`, userID)
		require.NoError(t, err)

		items, err := repo.ListInProgress(ctx, userID)
		require.NoError(t, err)
		assert.Empty(t, items)
	})
}
