CREATE TABLE IF NOT EXISTS watch_progress (
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  title_id INTEGER NOT NULL,
  position_sec INTEGER NOT NULL CHECK (position_sec >= 0),
  duration_sec INTEGER NOT NULL CHECK (duration_sec > 0),
  completed BOOLEAN NOT NULL DEFAULT FALSE,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, title_id),
  CHECK (position_sec <= duration_sec)
);

CREATE INDEX IF NOT EXISTS idx_watch_progress_user_in_progress ON watch_progress (user_id, updated_at DESC)
WHERE
  completed = FALSE;

