CREATE TABLE IF NOT EXISTS watch_progress (
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  title_id INTEGER NOT NULL,
  position_sec INTEGER NOT NULL CHECK (position_sec >= 0),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, title_id)
);

CREATE INDEX IF NOT EXISTS idx_watch_progress_user_updated ON watch_progress (user_id, updated_at DESC);
