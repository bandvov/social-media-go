-- Create notifications table
CREATE TABLE
    notifications (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL, -- The user receiving the notification
        type VARCHAR(50) NOT NULL CHECK (
            type IN (
                'new_follower',
                'mention',
                'direct_message',
                'post_comment',
                'comment_reply',
                'reaction'
            )
        ),
        entity_type VARCHAR(50) NOT NULL CHECK (
            entity_type IN ('user', 'post', 'comment', 'reaction')
        ),
        entity_id INT NOT NULL, -- ID of the related post, comment, reaction, or user
        message TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW (),
        is_read BOOLEAN DEFAULT FALSE
    );