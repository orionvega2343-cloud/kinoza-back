package domain

import "time"

type WatchlistItem struct {
	ID      int       `db:"id"`
	UserID  string    `db:"user_id"`
	TitleID int       `db:"title_id"`
	Status  string    `db:"status"`
	AddedAt time.Time `db:"added_at"`
}
