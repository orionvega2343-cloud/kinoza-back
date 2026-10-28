package domain

import "time"

type WatchProgress struct {
	UserID      string    `db:"user_id"`
	TitleID     int       `db:"title_id"`
	PositionSec int       `db:"position_sec"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (p *WatchProgress) Validate() error {
	if p.PositionSec < 0 {
		return ErrInvalidPosition
	}
	return nil
}
