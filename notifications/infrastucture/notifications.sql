-- Create notifications table
CREATE TABLE
    notifications (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL, -- The user receiving the notification
        type VARCHAR(50) NOT NULL CHECK (
            type IN (
                'new_follower',
                'new_mention',
                'new_direct_message',
                'new_post_comment',
                'new_comment_reply',
                'new_reaction'
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