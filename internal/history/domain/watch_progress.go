package domain

import "time"

const CompletedThreshold = 0.9

type WatchProgress struct {
	ID          int       `db:"id"`
	UserID      string    `db:"user_id"`
	TitleID     int       `db:"title_id"`
	PositionSec int       `db:"position_sec"`
	DurationSec int       `db:"duration_sec"`
	Completed   bool      `db:"completed"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (p *WatchProgress) Validate() error {
	if p.DurationSec <= 0 {
		return ErrInvalidDuration
	}

	if p.PositionSec < 0 || p.PositionSec > p.DurationSec {
		return ErrInvalidPosition
	}
	return nil
}

func (p *WatchProgress) MarkCompletion() {
	p.Completed = float64(p.PositionSec)/float64(p.DurationSec) >= CompletedThreshold
}
