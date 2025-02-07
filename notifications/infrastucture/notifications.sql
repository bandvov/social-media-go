-- Create notifications table

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL, -- The user receiving the notification
    type VARCHAR(50) NOT NULL CHECK (
    type IN (
        'new_follower',
        "new_reaction_like",
        "new_reaction_dislike",
        "new_reaction_love",
        "new_reaction_laugh",
        "new_reaction_angry",
        "new_reaction_wow",
        'new_direct_message',
        'new_post_comment',
        'new_comment_reply' 
        )),
         entity_type VARCHAR(50) NOT NULL CHECK (
        entity_type IN (
            'user',
            'post',
            'comment',
            'reaction'
        )), 
        entity_id INT NOT NULL, -- ID of the related post, comment, reaction, or user
        actor_ids INT[], -- Array of user IDs who triggered the event
        created_at TIMESTAMP DEFAULT NOW (),
        is_read BOOLEAN DEFAULT FALSE);