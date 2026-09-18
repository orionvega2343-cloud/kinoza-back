package domain

import (
	"time"
)

type WatchlistItem struct {
	ID      uint
	UserID  uint
	TitleID uint
	Status  string
	AddedAt time.Time
}
