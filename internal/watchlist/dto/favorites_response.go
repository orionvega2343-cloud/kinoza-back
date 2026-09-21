package dto

import "time"

type CreateWatchlistResponse struct {
	ID      int       `json:"id"`
	TitleID int       `json:"title_id"`
	Status  string    `json:"status"`
	AddedAt time.Time `json:"added_at"`
}
