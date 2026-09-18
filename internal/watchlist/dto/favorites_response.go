package dto

import (
	"time"
)

type CreateWatchlistResponse struct {
	ID      uint      `json:"id"`
	TitleID uint      `json:"title_id"`
	Status  string    `json:"status"`
	AddedAt time.Time `json:"added_at"`
}
