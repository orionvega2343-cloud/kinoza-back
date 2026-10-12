CREATE TABLE IF NOT EXISTS watchlist_items(
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title_id INTEGER NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('planned', 'watching', 'watched')),
    added_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, title_id)
);
