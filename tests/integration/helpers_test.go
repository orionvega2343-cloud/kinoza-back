//go:build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// resetDB чистит все таблицы, чтобы тесты не зависели друг от друга
func resetDB(t *testing.T) {
	t.Helper()

	_, err := db.Exec(`TRUNCATE users, refresh_tokens, watchlist_items, watch_progress RESTART IDENTITY CASCADE`)
	require.NoError(t, err)
}

func createUser(t *testing.T) string {
	t.Helper()

	var id string
	err := db.Get(
		&id,
		`INSERT INTO users(name, email, password_hash, role) VALUES($1, $2, $3, 'viewer') RETURNING id`,
		"test", uuid.NewString()+"@test.com", "hash",
	)
	require.NoError(t, err)
	return id
}
