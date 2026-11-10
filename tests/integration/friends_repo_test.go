//go:build integration

package integration

import (
	"context"
	"testing"

	"kinoza-back/internal/friends/domain"
	"kinoza-back/internal/friends/repository"
	"kinoza-back/pkg/transaction"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const unknownUserID = "00000000-0000-0000-0000-000000000000"

func TestFriendRequestRepo_Create(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	repo := repository.NewFriendRequestRepository(db)
	alice, bob, carol := createUser(t), createUser(t), createUser(t)

	t.Run("creates pending request", func(t *testing.T) {
		req := &domain.FriendRequest{FromUserID: alice, ToUserID: bob}

		require.NoError(t, repo.Create(ctx, req))
		assert.NotZero(t, req.ID)
		assert.Equal(t, domain.StatusPending, req.Status)
		assert.False(t, req.CreatedAt.IsZero())
	})

	t.Run("duplicate pending returns ErrRequestAlreadyExists", func(t *testing.T) {
		err := repo.Create(ctx, &domain.FriendRequest{FromUserID: alice, ToUserID: bob})

		assert.ErrorIs(t, err, domain.ErrRequestAlreadyExists)
	})

	t.Run("reverse pending returns ErrRequestAlreadyExists", func(t *testing.T) {
		err := repo.Create(ctx, &domain.FriendRequest{FromUserID: bob, ToUserID: alice})

		assert.ErrorIs(t, err, domain.ErrRequestAlreadyExists)
	})

	t.Run("new request allowed after decline", func(t *testing.T) {
		req := &domain.FriendRequest{FromUserID: alice, ToUserID: carol}
		require.NoError(t, repo.Create(ctx, req))
		req.Status = domain.StatusDeclined
		require.NoError(t, repo.UpdateStatus(ctx, req))

		assert.NoError(t, repo.Create(ctx, &domain.FriendRequest{FromUserID: alice, ToUserID: carol}))
	})

	t.Run("unknown recipient returns ErrUserNotFound", func(t *testing.T) {
		err := repo.Create(ctx, &domain.FriendRequest{FromUserID: alice, ToUserID: unknownUserID})

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("self request rejected by CHECK", func(t *testing.T) {
		assert.Error(t, repo.Create(ctx, &domain.FriendRequest{FromUserID: alice, ToUserID: alice}))
	})
}

func TestFriendRequestRepo_GetByIDForUpdate(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	repo := repository.NewFriendRequestRepository(db)
	tx := transaction.NewTransactor(db)
	alice, bob := createUser(t), createUser(t)
	created := &domain.FriendRequest{FromUserID: alice, ToUserID: bob}
	require.NoError(t, repo.Create(ctx, created))

	t.Run("returns request inside transaction", func(t *testing.T) {
		var got *domain.FriendRequest
		err := tx.Transaction(ctx, func(ctx context.Context) error {
			var err error
			got, err = repo.GetByIDForUpdate(ctx, created.ID)
			return err
		})

		require.NoError(t, err)
		assert.Equal(t, created.ID, got.ID)
		assert.Equal(t, alice, got.FromUserID)
		assert.Equal(t, bob, got.ToUserID)
		assert.Equal(t, domain.StatusPending, got.Status)
	})

	t.Run("missing returns ErrRequestNotFound", func(t *testing.T) {
		_, err := repo.GetByIDForUpdate(ctx, 999)

		assert.ErrorIs(t, err, domain.ErrRequestNotFound)
	})
}

func TestFriendRequestRepo_UpdateStatus(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	repo := repository.NewFriendRequestRepository(db)
	alice, bob := createUser(t), createUser(t)
	req := &domain.FriendRequest{FromUserID: alice, ToUserID: bob}
	require.NoError(t, repo.Create(ctx, req))

	t.Run("updates status", func(t *testing.T) {
		req.Status = domain.StatusAccepted
		require.NoError(t, repo.UpdateStatus(ctx, req))

		var status string
		require.NoError(t, db.Get(&status, `SELECT status FROM friend_requests WHERE id = $1`, req.ID))
		assert.Equal(t, domain.StatusAccepted, status)
	})

	t.Run("missing returns ErrRequestNotFound", func(t *testing.T) {
		err := repo.UpdateStatus(ctx, &domain.FriendRequest{ID: 999, Status: domain.StatusAccepted})

		assert.ErrorIs(t, err, domain.ErrRequestNotFound)
	})
}

func TestFriendRequestRepo_Lists(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	repo := repository.NewFriendRequestRepository(db)
	alice, bob, carol, dave := createUser(t), createUser(t), createUser(t), createUser(t)

	require.NoError(t, repo.Create(ctx, &domain.FriendRequest{FromUserID: bob, ToUserID: alice}))
	require.NoError(t, repo.Create(ctx, &domain.FriendRequest{FromUserID: carol, ToUserID: alice}))
	require.NoError(t, repo.Create(ctx, &domain.FriendRequest{FromUserID: alice, ToUserID: dave}))

	declined := &domain.FriendRequest{FromUserID: dave, ToUserID: bob}
	require.NoError(t, repo.Create(ctx, declined))
	declined.Status = domain.StatusDeclined
	require.NoError(t, repo.UpdateStatus(ctx, declined))

	t.Run("incoming returns only pending to user", func(t *testing.T) {
		got, err := repo.ListIncoming(ctx, alice)

		require.NoError(t, err)
		require.Len(t, got, 2)
		for _, r := range got {
			assert.Equal(t, alice, r.ToUserID)
		}
	})

	t.Run("outgoing returns only pending from user", func(t *testing.T) {
		got, err := repo.ListOutgoing(ctx, alice)

		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, dave, got[0].ToUserID)
	})

	t.Run("declined not listed", func(t *testing.T) {
		got, err := repo.ListIncoming(ctx, bob)

		require.NoError(t, err)
		assert.Empty(t, got)
		assert.NotNil(t, got)
	})
}

func TestFriendshipRepo_CreateExists(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	repo := repository.NewFriendshipRepository(db)
	alice, bob := createUser(t), createUser(t)

	t.Run("not exists before create", func(t *testing.T) {
		ok, err := repo.Exists(ctx, domain.NewFriendship(alice, bob))

		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("creates friendship", func(t *testing.T) {
		f := domain.NewFriendship(alice, bob)

		require.NoError(t, repo.Create(ctx, f))
		assert.NotZero(t, f.ID)
		assert.False(t, f.CreatedAt.IsZero())
	})

	t.Run("exists in both argument orders", func(t *testing.T) {
		ok1, err := repo.Exists(ctx, domain.NewFriendship(alice, bob))
		require.NoError(t, err)
		ok2, err := repo.Exists(ctx, domain.NewFriendship(bob, alice))
		require.NoError(t, err)

		assert.True(t, ok1)
		assert.True(t, ok2)
	})

	t.Run("duplicate returns ErrAlreadyFriends", func(t *testing.T) {
		err := repo.Create(ctx, domain.NewFriendship(bob, alice))

		assert.ErrorIs(t, err, domain.ErrAlreadyFriends)
	})
}

func TestFriendshipRepo_ListDelete(t *testing.T) {
	resetDB(t)
	ctx := context.Background()
	repo := repository.NewFriendshipRepository(db)
	alice, bob, carol := createUser(t), createUser(t), createUser(t)
	require.NoError(t, repo.Create(ctx, domain.NewFriendship(alice, bob)))
	require.NoError(t, repo.Create(ctx, domain.NewFriendship(carol, alice)))

	t.Run("list returns friendships on both sides", func(t *testing.T) {
		got, err := repo.List(ctx, alice)

		require.NoError(t, err)
		require.Len(t, got, 2)
		friends := []string{got[0].FriendOf(alice), got[1].FriendOf(alice)}
		assert.ElementsMatch(t, []string{bob, carol}, friends)
	})

	t.Run("delete removes friendship", func(t *testing.T) {
		require.NoError(t, repo.Delete(ctx, domain.NewFriendship(bob, alice)))

		ok, err := repo.Exists(ctx, domain.NewFriendship(alice, bob))
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("delete missing returns ErrFriendshipNotFound", func(t *testing.T) {
		err := repo.Delete(ctx, domain.NewFriendship(alice, bob))

		assert.ErrorIs(t, err, domain.ErrFriendshipNotFound)
	})
}
