CREATE TABLE
    IF NOT EXISTS user_sessions (
        id bigserial PRIMARY KEY NOT NULL,
        token_id UUID NOT NULL,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        last_used_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        expired_at TIMESTAMPTZ NOT NULL
    );

CREATE UNIQUE INDEX idx_token_id_user_id ON user_sessions(token_id, user_id)
