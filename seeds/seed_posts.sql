INSERT INTO posts (author_id, content, visibility, tags, like_count, comment_count, share_count, pinned)
VALUES 
(1, 'Enjoying the sunny day at the beach!', 'public', '#beach,#sunnyday', 10, 5, 3, FALSE),
(2, 'Just tried a new recipe, and it is amazing!', 'public', '#cooking,#recipe', 20, 8, 2, FALSE),
(1, 'Poll: Which do you prefer? Coffee or Tea?', 'public', '#poll,#coffee,#tea', 15, 12, 6, FALSE),
(1, 'Throwback to last year’s hiking trip.', 'private', '#hiking,#nature', 5, 2, 1, FALSE),
(2, 'Scheduled post for later!', 'public', '#scheduled',  0, 0, 0, TRUE);
