package repository

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"kinoza-back/internal/watchlist/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

var errDB = errors.New("db error")

func setupMock(t *testing.T) (domain.WatchlistRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewWatchlistRepository(sqlx.NewDb(db, "sqlmock")), mock
}

func TestRepoCreate(t *testing.T) {
	tests := []struct {
		name    string
		rows    *sqlmock.Rows
		dbErr   error
		wantErr error
		wantID  int
	}{
		{name: "success", rows: sqlmock.NewRows([]string{"id", "added_at"}).AddRow(7, time.Now()), wantErr: nil, wantID: 7},
		{name: "db error", dbErr: errDB, wantErr: errDB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock := setupMock(t)

			q := mock.ExpectQuery(`INSERT INTO watchlist_items(.+)ON CONFLICT \(user_id, title_id\) DO UPDATE SET status = EXCLUDED.status, updated_at = NOW\(\)`).
				WithArgs("user-1", 42, domain.StatusPlanned)
			if tt.dbErr != nil {
				q.WillReturnError(tt.dbErr)
			} else {
				q.WillReturnRows(tt.rows)
			}

			item := &domain.WatchlistItem{UserID: "user-1", TitleID: 42, Status: domain.StatusPlanned}
			err := repo.Create(context.Background(), item)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if item.ID != tt.wantID {
				t.Errorf("expected id %d, got %d", tt.wantID, item.ID)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestRepoUpdate(t *testing.T) {
	tests := []struct {
		name    string
		result  driver.Result
		dbErr   error
		wantErr error
	}{
		{name: "success", result: sqlmock.NewResult(0, 1), wantErr: nil},
		{name: "not found or other user", result: sqlmock.NewResult(0, 0), wantErr: domain.ErrNotFound},
		{name: "db error", dbErr: errDB, wantErr: errDB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock := setupMock(t)

			e := mock.ExpectExec(`UPDATE watchlist_items SET status = \$1, updated_at = NOW\(\) WHERE user_id = \$2 AND title_id = \$3`).
				WithArgs(domain.StatusWatched, "user-1", 42)
			if tt.dbErr != nil {
				e.WillReturnError(tt.dbErr)
			} else {
				e.WillReturnResult(tt.result)
			}

			err := repo.Update(context.Background(), &domain.WatchlistItem{TitleID: 42, UserID: "user-1", Status: domain.StatusWatched})

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestRepoDelete(t *testing.T) {
	tests := []struct {
		name    string
		result  driver.Result
		dbErr   error
		wantErr error
	}{
		{name: "success", result: sqlmock.NewResult(0, 1), wantErr: nil},
		{name: "not found or other user", result: sqlmock.NewResult(0, 0), wantErr: domain.ErrNotFound},
		{name: "db error", dbErr: errDB, wantErr: errDB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock := setupMock(t)

			e := mock.ExpectExec(`DELETE FROM watchlist_items WHERE user_id = \$1 AND title_id = \$2`).
				WithArgs("user-1", 42)
			if tt.dbErr != nil {
				e.WillReturnError(tt.dbErr)
			} else {
				e.WillReturnResult(tt.result)
			}

			err := repo.Delete(context.Background(), 42, "user-1")

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestRepoList(t *testing.T) {
	columns := []string{"id", "user_id", "title_id", "status", "added_at"}

	tests := []struct {
		name      string
		rows      *sqlmock.Rows
		dbErr     error
		wantErr   error
		wantCount int
	}{
		{name: "two items", rows: sqlmock.NewRows(columns).AddRow(2, "user-1", 20, "watched", time.Now()).AddRow(1, "user-1", 10, "planned", time.Now()), wantCount: 2},
		{name: "empty", rows: sqlmock.NewRows(columns), wantCount: 0},
		{name: "db error", dbErr: errDB, wantErr: errDB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mock := setupMock(t)

			q := mock.ExpectQuery(`SELECT (.+) FROM watchlist_items WHERE user_id = \$1 ORDER BY added_at DESC`).
				WithArgs("user-1")
			if tt.dbErr != nil {
				q.WillReturnError(tt.dbErr)
			} else {
				q.WillReturnRows(tt.rows)
			}

			items, err := repo.List(context.Background(), "user-1")

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if len(items) != tt.wantCount {
				t.Errorf("expected %d items, got %d", tt.wantCount, len(items))
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}
