
CREATE TABLE
    IF NOT EXISTS tweets (
        id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v1(),
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        body VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );