CREATE TABLE watchlist_items(
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    title_id INTEGER NOT NULL,
    status TEXT NOT NULL,
    added_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, title_id)
);
