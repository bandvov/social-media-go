CREATE TABLE IF NOT EXISTS public.followers
    ( 
    follower_id INT, 
    followee_id INT,
     PRIMARY KEY (follower_id, followee_id),
     FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
     FOREIGN KEY (followee_id) REFERENCES users(id) ON DELETE CASCADE
    );