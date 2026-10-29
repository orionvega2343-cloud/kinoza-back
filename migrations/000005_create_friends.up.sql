CREATE TABLE IF NOT EXISTS friend_requests (
    id SERIAL PRIMARY KEY,
    from_user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    to_user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined')),
     created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
     CHECK (from_user_id <> to_user_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS  uniq_friend_requests_pending_pair
    ON friend_requests (LEAST(from_user_id, to_user_id), GREATEST(from_user_id, to_user_id))
    WHERE status = 'pending';

    CREATE INDEX IF NOT EXISTS idx_friend_requests_to_user ON friend_requests (to_user_id, status);
    CREATE INDEX IF NOT EXISTS idx_friend_requests_from_user ON friend_requests (from_user_id, status);

    CREATE TABLE IF NOT EXISTS friendships (
        id SERIAL PRIMARY KEY,
        user_id_a UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        user_id_b UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CHECK (user_id_a < user_id_b),
        UNIQUE (user_id_a, user_id_b)
    );

    CREATE INDEX IF NOT EXISTS idx_friendships_users_b ON friendships (user_id_b);
