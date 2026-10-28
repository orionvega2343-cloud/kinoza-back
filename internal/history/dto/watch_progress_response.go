package dto

import (
	"time"
)

type WatchProgressResponse struct {
	TitleID     int       `json:"title_id"`
	PositionSec int       `json:"position_sec"`
	UpdatedAt   time.Time `json:"updated_at"`
}
