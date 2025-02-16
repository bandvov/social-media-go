CREATE TABLE
    activities (
        id SERIAL PRIMARY KEY,
        user_id TEXT NOT NULL,
        action TEXT NOT NULL,
        target_id TEXT NOT NULL,
        target_type TEXT NOT NULL CHECK (
            target_type IN ('post', 'comment', 'user', 'reaction')
        ),
        event_data JSONB,
        created_at TIMESTAMP NOT NULL DEFAULT NOW ()
    );